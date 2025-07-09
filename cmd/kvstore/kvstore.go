package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go_util/pkg/relay"
	"go_util/pkg/rlog"
	"go_util/pkg/rlog_util/kvstore"
	"go_util/pkg/rlog_util/transport"
	"go_util/pkg/uuid"
	"gopkg.in/yaml.v2"
)

type config struct {
	Relay             string `json:"relay"`
	RpcTimeoutMs      int    `json:"rpc_timeout_ms"`
	CompactionBlock   int    `json:"compaction_block"`
	CompactionRatio   int    `json:"compaction_ratio"`
	RetryUntilUpdate  int    `json:"retry_until_update"`
	MinRetryTimeoutMs int    `json:"min_retry_timeout_ms"`
	MaxRetryTimeoutMs int    `json:"max_retry_timeout_ms"`
	UpdateIntervalMs  int    `json:"update_interval_ms"`
	RestartTimeoutMs  int    `json:"restart_timeout_ms"`
	Cluster           []struct {
		Name string `json:"name"`
		Port int    `json:"port"`
	} `json:"cluster"`
}

func fetchIdAndConfig() (int, config) {
	i := flag.Int("id", 0, "id of node")
	c := flag.String("config", "test_config.json", "config file")
	flag.Parse()
	b, err := ioutil.ReadFile(*c)
	if err != nil {
		panic(err)
	}
	var cfg config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}
	return *i, cfg
}

type container struct {
	nodeId    rlog.NodeId
	config    config
	store     kvstore.KVStore
	transport transport.Transport
	node      *rlog.Node
	writeMu   *sync.Mutex
}

func newContainer(id int, cfg config) (container, error) {
	// no di
	nodeId := rlog.NodeId(id)
	addressBook := map[rlog.NodeId]string{}
	for peerId := range cfg.Cluster {
		addressBook[rlog.NodeId(peerId)] = cfg.Cluster[peerId].Name
	}
	agent, err := relay.NewPeer(cfg.Cluster[nodeId].Name, "", cfg.Relay)
	if err != nil {
		return container{}, err
	}
	trans := transport.NewRelay(agent)
	store := kvstore.New()
	node := &rlog.Node{
		NodeId: nodeId,
		Router: trans.Router,
		Cluster: rlog.ClusterState{
			AddressBook:      addressBook,
			RpcTimeout:       time.Duration(cfg.RpcTimeoutMs) * time.Millisecond,
			CompactionBlock:  rlog.LogId(cfg.CompactionBlock),
			CompactionRatio:  rlog.LogId(cfg.CompactionRatio),
			RetryUntilUpdate: cfg.RetryUntilUpdate,
		},
		Acceptor: rlog.AcceptorState{
			Value: rlog.Value{
				Object: store,
			},
		},
		Proposer: rlog.ProposerState{
			AcceptIdMap: map[rlog.NodeId]rlog.LogId{},
		},
	}
	return container{
		nodeId:    nodeId,
		config:    cfg,
		store:     store,
		transport: trans,
		node:      node,
		writeMu:   &sync.Mutex{},
	}, nil
}

func (ctn container) run() {
	// transport
	go func() {
		for {
			err := ctn.transport.ListenAndServe(ctn.node.Handler)
			if err != nil {
				fmt.Printf("[transport] error: %v\n", err)
			}
			time.Sleep(time.Duration(ctn.config.RestartTimeoutMs) * time.Millisecond)
			fmt.Printf("[transport] restarted\n")
		}
	}()
	// http
	go func() {
		for {
			err := newServer(ctn).ListenAndServe()
			if err != nil {
				fmt.Printf("[application] error: %v\n", err)
			}
			time.Sleep(time.Duration(ctn.config.RestartTimeoutMs) * time.Millisecond)
			fmt.Printf("[application] restarted\n")
		}
	}()
	// update loop
	go func() {
		for {
			updateCtx, updateCancel := context.WithTimeout(context.Background(), time.Duration(ctn.config.UpdateIntervalMs)*time.Millisecond)
			ctn.node.Update(updateCtx)
			updateCancel()
			time.Sleep(time.Duration(ctn.config.UpdateIntervalMs) * time.Millisecond)
		}
	}()
	<-context.Background().Done()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	id, cfg := fetchIdAndConfig()
	ctn, err := newContainer(id, cfg)
	if err != nil {
		fmt.Printf("[main] stopped with error %v\n", err)
		panic(err)
	}
	ctn.run()
}

func proposeAndWrite(ctn container, ctx *gin.Context, cmd kvstore.Command) {
	wrappedCommand := cmd.Encode()
	ctn.writeMu.Lock()
	defer ctn.writeMu.Unlock()
	ch, cancel := ctn.store.Watch(cmd.Uuid)
	defer cancel()
	err := ctn.node.Propose(ctx.Request.Context(), wrappedCommand, exponentialBackoff(
		time.Duration(ctn.config.MinRetryTimeoutMs)*time.Millisecond,
		time.Duration(ctn.config.MaxRetryTimeoutMs)*time.Millisecond,
		2.0,
	))
	if err != nil {
		ctx.String(http.StatusServiceUnavailable, err.Error())
		return
	}
	consistent := <-ch
	if !consistent {
		ctx.String(http.StatusConflict, "")
		return
	}
	ctx.String(http.StatusOK, "")
}

func newServer(ctn container) *http.Server {
	e := gin.New()
	e.Handle(http.MethodPost, "/kvstore/*key", func(c *gin.Context) {
		defer c.Request.Body.Close()
		key := c.Param("key")
		versionStr := c.Query("version")
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		val := string(b)
		proposeAndWrite(ctn, c, kvstore.Command{
			Uuid:      uuid.New(),
			Operation: kvstore.OpSet,
			Key:       key,
			Version:   uint64(version),
			Value:     val,
		})
	})
	e.Handle(http.MethodDelete, "/kvstore/*key", func(c *gin.Context) {
		defer c.Request.Body.Close()
		key := c.Param("key")
		versionStr := c.Query("version")
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		proposeAndWrite(ctn, c, kvstore.Command{
			Uuid:      uuid.New(),
			Operation: kvstore.OpDel,
			Key:       key,
			Version:   uint64(version),
		})
	})
	e.Handle(http.MethodGet, "/kvstore", func(c *gin.Context) {
		defer c.Request.Body.Close()
		o := ctn.store.Snapshot()
		// human readable
		b, err := yaml.Marshal(o)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, string(b))
	})
	e.Handle(http.MethodGet, "/kvstore/*key", func(c *gin.Context) {
		key := c.Param("key")
		o := ctn.store.Get(key)
		// json
		c.JSON(http.StatusOK, o)
	})
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(ctn.config.Cluster[ctn.nodeId].Port),
		Handler: e,
	}
	return server
}

func exponentialBackoff(minTimeout time.Duration, maxTimeout time.Duration, scale float64) rlog.RetryPolicy {
	if minTimeout == 0 || maxTimeout == 0 {
		panic("min timeout and max timeout must be positive")
	}
	timeout := minTimeout
	return func() time.Duration {
		duration := time.Duration(rand.Intn(int(timeout)))
		timeout = time.Duration(float64(timeout) * scale)
		if timeout > maxTimeout {
			timeout = maxTimeout
		}
		return duration
	}
}
