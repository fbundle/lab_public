package proto

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type message struct {
	Type    Type        `json:"type"`
	Payload interface{} `json:"payload"`
}

type Type string

type Proto interface {
	MustRegister(mType Type, mPayload interface{})
	Marshal(m interface{}) (b []byte, err error)
	Unmarshal(b []byte) (m interface{}, err error)
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

func NewProto() Proto {
	return &proto{
		protoMap: make(map[Type]interface{}),
	}
}

type proto struct {
	protoMap map[Type]interface{}
}

func (p *proto) MustRegister(mType Type, mPayload interface{}) {
	if err := mustBePtrOfStruct(mPayload); err != nil {
		panic(err)
	}
	if ptr, ok := p.protoMap[mType]; ok {
		if reflect.TypeOf(mPayload) != reflect.TypeOf(ptr) {
			panic(fmt.Errorf("type exists %s", mType))
		}
	} else {
		p.protoMap[mType] = mPayload
	}
}

func (p *proto) Marshal(payload interface{}) (b []byte, err error) {
	if err := mustBePtrOfStruct(payload); err != nil {
		panic(err)
	}
	for mType, mPayload := range p.protoMap {
		if reflect.TypeOf(mPayload) == reflect.TypeOf(payload) {
			m := &message{
				Type:    mType,
				Payload: payload,
			}
			return json.Marshal(m)
		}
	}
	return nil, fmt.Errorf("type was not registered")
}

func (p *proto) Unmarshal(b []byte) (payload interface{}, err error) {
	h := &message{
		Type: "",
	}
	err = json.Unmarshal(b, h)
	if err != nil {
		return nil, err
	}
	if _, ok := p.protoMap[h.Type]; !ok {
		return nil, fmt.Errorf("unknown type %s\n", string(b))
	}
	m := &message{
		Payload: reflect.New(reflect.TypeOf(p.protoMap[h.Type]).Elem()).Interface(),
	}
	err = json.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	return m.Payload, nil
}
