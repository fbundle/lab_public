package main

import (
	"bytes"
	"ca/pkg/kvstore_client"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type config struct {
	KVStore  string            `json:"kvstore"`
	Once     map[string]string `json:"once"`
	Interval map[string]string `json:"interval"`
}

var c *string

func init() {
	c = flag.String("config", "test_config.json", "config file")
	flag.Parse()
}

func fetchConfig() config {
	b, err := ioutil.ReadFile(*c)
	if err != nil {
		panic(err)
	}
	var cfg config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func propose(client kvstore_client.Client, key string, command string) {
	getCtx, getCancel := context.WithTimeout(context.Background(), 5*time.Second)
	version, _, err := client.Get(getCtx, key)
	getCancel()
	if err != nil {
		fmt.Println(key, command, err)
		return
	}
	var value string
	buffer := &bytes.Buffer{}
	fields := strings.Fields(command)
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Stdout = buffer
	err = cmd.Run()
	if err != nil {
		fmt.Println(key, command, err)
		value = err.Error()
	} else {
		value = buffer.String()
	}
	setCtx, setCancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Set(setCtx, key, version+1, value)
	setCancel()
	if err != nil {
		fmt.Println(key, command, err)
		return
	}
	fmt.Println("proposed", key, command)
}

func main() {
	cfg := fetchConfig()
	client := kvstore_client.New(cfg.KVStore)
	for key, command := range cfg.Once {
		propose(client, key, command)
	}
	for {
		for key, command := range cfg.Interval {
			propose(client, key, command)
		}
		time.Sleep(10 * time.Second)
	}

}
