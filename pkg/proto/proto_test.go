package proto_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fbundle/go_util/pkg/proto"
)

type structA struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type structB struct {
	Value float64 `json:"value"`
}

type header struct {
	Sender string `json:"sender"`
}

func TestProtoOk(t *testing.T) {
	p := proto.NewProto()
	p.MustRegister("struct_a", &structA{})
	expected := &structA{
		Name: "hello",
		Age:  12,
	}
	b, err := p.Marshal(expected)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))
	i, err := p.Unmarshal(b)
	if err != nil {
		t.Error(err)
		return
	}
	actual := i.(*structA)
	if !reflect.DeepEqual(expected, actual) {
		t.Error("wrong")
		return
	}
}

func TestProtoErr(t *testing.T) {
	p := proto.NewProto()
	_, err := p.Marshal(&structB{Value: 1.3})
	if err == nil {
		t.Error("error expected")
		return
	}
}
