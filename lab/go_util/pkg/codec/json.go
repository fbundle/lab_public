package codec

import "encoding/json"

func NewJsonCodec() Codec {
	return &jsonCodec{}
}

type jsonCodec struct{}

func (c *jsonCodec) Marshal(o interface{}) (b []byte, err error) {
	return json.Marshal(o)
}

func (c *jsonCodec) Unmarshal(b []byte, o interface{}) (err error) {
	return json.Unmarshal(b, o)
}
