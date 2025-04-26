package transport

import (
	"bytes"
	"ca/pkg/rlog/rpc"
	rlog_codec "ca/pkg/rlog_util/codec"
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	protoScheme     = "http://"
	transportPath   = "/transport"
	transportMethod = http.MethodPost
)

func NewHttp(host rpc.Address) Transport {
	return &httpTransport{
		host:   host,
		server: nil,
	}
}

type httpTransport struct {
	host   rpc.Address
	server *http.Server
}

func (s *httpTransport) Close() error {
	return s.server.Close()
}

func (s *httpTransport) ListenAndServe(handler rpc.Handler) error {
	e := gin.New()
	e.Handle(transportMethod, transportPath, func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		b, err = base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		request, err := rlog_codec.Unmarshal(b)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		call := &rpc.Call{
			Request: request,
		}
		rpc.WaitThenCancel(call, c.Request.Context())
		handler(call)
		<-call.Done()
		b, err = rlog_codec.Marshal(call.Response)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		s := base64.StdEncoding.EncodeToString(b)
		c.String(http.StatusOK, s)
	})

	s.server = &http.Server{
		Addr:    getListenAddress(s.host),
		Handler: e,
	}
	return s.server.ListenAndServe()
}

func (s *httpTransport) Router(receiver rpc.Address) rpc.Handler {
	return func(call *rpc.Call) {
		response := func() interface{} {
			b, err := rlog_codec.Marshal(call.Request)
			if err != nil {
				return nil
			}
			s := base64.StdEncoding.EncodeToString(b)
			req, err := http.NewRequest(transportMethod, protoScheme+receiver+transportPath, bytes.NewBuffer([]byte(s)))
			if err != nil {
				return nil
			}
			client := &http.Client{}
			ctx, cancel := context.WithCancel(context.Background())
			req = req.WithContext(ctx)
			go func() {
				<-call.Done()
				cancel()
			}()
			res, err := client.Do(req)
			if err != nil {
				return nil
			}
			if res.StatusCode != http.StatusOK {
				log.Println(res.StatusCode)
				return nil
			}
			b, err = ioutil.ReadAll(res.Body)
			if err != nil {
				return nil
			}
			defer res.Body.Close()
			b, err = base64.StdEncoding.DecodeString(string(b))
			if err != nil {
				return nil
			}
			response, err := rlog_codec.Unmarshal(b)
			if err != nil {
				response = nil
				return nil
			}
			return response
		}()
		call.Write(response)
	}
}
func getListenAddress(address string) string {
	return ":" + strings.SplitN(address, ":", 2)[1]
}
