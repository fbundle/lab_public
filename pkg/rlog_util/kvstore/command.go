package kvstore

import (
	"encoding/json"
)

type Operation string

const (
	OpSet Operation = "set"
	OpDel Operation = "del"
)

type Command struct {
	Uuid      string    `json:"uuid"`
	Operation Operation `json:"operation"`
	Key       string    `json:"key"`
	Version   uint64    `json:"version"`
	Value     string    `json:"value"`
}

func (c Command) Encode() string {
	b, _ := json.Marshal(c)
	return string(b)
}
func (c *Command) Decode(s string) error {
	return json.Unmarshal([]byte(s), c)
}
