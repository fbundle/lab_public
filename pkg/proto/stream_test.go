package proto_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/fbundle/go_util/pkg/proto"
)

type Message struct {
	Data string `json:"data"`
}

var protoStream proto.Proto

func init() {
	protoStream = proto.NewProto()
	protoStream.MustRegister("message", &Message{})
}

func TestStream(t *testing.T) {
	var b = bytes.NewBuffer([]byte{})

	e := protoStream.NewEncoder(b)
	d := protoStream.NewDecoder(b)

	_ = e.Encode(&Message{Data: "hello"})
	_ = e.Encode(&Message{Data: "this is khanh"})
	for {
		m, err := d.Decode()
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(m)
	}
}
