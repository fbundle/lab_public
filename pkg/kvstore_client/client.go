package kvstore_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/khanh-nguyen-code/go_util/pkg/rlog_util/kvstore"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client interface {
	Get(ctx context.Context, key string) (version uint64, value string, err error)
	Set(ctx context.Context, key string, version uint64, value string) error
	Del(ctx context.Context, key string, version uint64) error
}

func New(path string) Client {
	return &client{path: path}
}

type client struct {
	path string
}

func (c *client) Get(ctx context.Context, path string) (version uint64, value string, err error) {
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}
	req, err := http.NewRequest(
		http.MethodGet,
		c.path+path,
		nil,
	)
	if err != nil {
		return 0, "", err
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status_code_%d", res.StatusCode)
		return 0, "", err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}
	e := kvstore.Entry{}
	err = json.Unmarshal(b, &e)
	return e.Version, e.Value, err
}
func (c *client) Set(ctx context.Context, path string, version uint64, value string) error {
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}
	req, err := http.NewRequest(
		http.MethodPost,
		c.path+path+"?version="+strconv.Itoa(int(version)),
		bytes.NewBuffer([]byte(value)),
	)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status_code_%d", res.StatusCode)
	}
	return nil
}
func (c *client) Del(ctx context.Context, key string, version uint64) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		c.path+"/"+key+"?version="+strconv.Itoa(int(version)),
		nil,
	)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status_code_%d", res.StatusCode)
	}
	return nil
}
