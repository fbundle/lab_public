package codec

import (
	"encoding/xml"
)

func NewXmlCodec() Codec {
	return &xmlCodec{}
}

type xmlCodec struct{}

func (c *xmlCodec) Marshal(o interface{}) (b []byte, err error) {
	return xml.Marshal(o)
}

func (c *xmlCodec) Unmarshal(b []byte, o interface{}) (err error) {
	return xml.Unmarshal(b, o)
}
