package rlog_test

import (
	"context"
	"encoding/json"
	"fmt"
	rlog "github.com/khanh-nguyen-code/go_util/pkg/rlog"
	"github.com/khanh-nguyen-code/go_util/pkg/rlog/rpc"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestNode_ProposeOnce(t *testing.T) {
	rand.Seed(1234)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	numNodes := 3
	cluster := makeCluster(numNodes, transportModel{
		drop:     0.5,
		minDelay: 0,
		maxDelay: 10 * time.Millisecond,
	})
	nodeMap := make(map[rlog.NodeId]*rlog.Node)
	for nodeId := range cluster.addressBook {
		nodeId := nodeId
		node := &rlog.Node{
			NodeId: nodeId,
			Router: cluster.router,
			Cluster: rlog.ClusterState{
				AddressBook:      cluster.addressBook,
				RpcTimeout:       10 * time.Millisecond,
				CompactionBlock:  2,
				CompactionRatio:  2,
				RetryUntilUpdate: 2,
			},
			Acceptor: rlog.AcceptorState{
				Value: rlog.Value{
					Object: newLog(),
				},
			},
			Proposer: rlog.ProposerState{
				AcceptIdMap: map[rlog.NodeId]rlog.LogId{},
			},
		}
		nodeMap[nodeId] = node
		go cluster.serve(ctx, nodeId, node.Handler)
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					node.Update(ctx)
				}
			}
		}()
	}

	doneCount := uint32(0)
	for nodeId := range cluster.addressBook {
		go func(nodeId rlog.NodeId) {
			for i := 0; i < 3; i++ {
				v := fmt.Sprintf("command_%d_%d", nodeId, i)
				_ = nodeMap[nodeId].Propose(context.Background(), v, exponentialBackoff(10*time.Millisecond, 80*time.Millisecond, 2))
			}

			if atomic.AddUint32(&doneCount, 1) >= uint32(len(cluster.addressBook)) {
				cancel()
			}
		}(nodeId)
	}
	<-ctx.Done()

	for _, replicatedLog := range nodeMap {
		b, err := json.MarshalIndent(replicatedLog.Acceptor, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
}

type cluster struct {
	model       transportModel
	addressBook map[rlog.NodeId]rpc.Address
	chanMap     map[rlog.NodeId]chan *rpc.Call
}

func makeCluster(numNodes int, model transportModel) *cluster {
	c := &cluster{
		model:       model,
		addressBook: map[rlog.NodeId]rpc.Address{},
		chanMap:     map[rlog.NodeId]chan *rpc.Call{},
	}
	for i := 0; i < numNodes; i++ {
		nodeId := rlog.NodeId(i)
		c.addressBook[nodeId] = strconv.Itoa(int(nodeId))
		c.chanMap[nodeId] = make(chan *rpc.Call, 1024)
	}
	return c
}

func (c *cluster) router(receiverAddr rpc.Address) rpc.Handler {
	receiver, _ := strconv.Atoi(receiverAddr)
	return func(rpc *rpc.Call) {
		c.model.do(func() {
			c.chanMap[rlog.NodeId(receiver)] <- rpc
		}, func() {
			rpc.Write(nil)
		})
	}
}

func (c *cluster) serve(ctx context.Context, nodeId rlog.NodeId, handler rpc.Handler) {
	for {
		select {
		case <-ctx.Done():
			return
		case call := <-c.chanMap[nodeId]:
			handler(call)
		}
	}
}

type transportModel struct {
	drop     float64
	minDelay time.Duration
	maxDelay time.Duration
}

func (tm transportModel) do(sendCb func(), dropCb func()) {
	time.Sleep(tm.minDelay + time.Duration(rand.Intn(int(tm.maxDelay-tm.minDelay))))
	if rand.Float64() < tm.drop {
		dropCb()
	} else {
		sendCb()
	}
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

type log struct {
	Entries []string
}

func (l log) Marshal() ([]byte, error) {
	return json.Marshal(l.Entries)
}

func (l *log) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, &l.Entries)
}

func (l *log) Commit(command ...rlog.Command) {
	l.Entries = append(l.Entries, command...)
}
func newLog() rlog.Object {
	return &log{}
}
