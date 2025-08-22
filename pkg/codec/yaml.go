package codec

import (
	"github.com/go-yaml/yaml"
)

func NewYamlCodec() Codec {
	return &yamlCodec{}
}

type yamlCodec struct{}

func (c *yamlCodec) Marshal(o interface{}) (b []byte, err error) {
	return yaml.Marshal(o)
}

func (c *yamlCodec) Unmarshal(b []byte, o interface{}) (err error) {
	return yaml.Unmarshal(b, o)
}
