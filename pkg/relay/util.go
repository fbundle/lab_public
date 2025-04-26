package relay

import (
	"ca/pkg/relay/proto/gen/relay_pb"
	"encoding/binary"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
)

var sizeError = errors.New("size")

func readExact(reader io.Reader, size int) ([]byte, error) {
	buffer := make([]byte, size)
	n, err := io.ReadFull(reader, buffer)
	return buffer[:n], err
}

// readAndUnmarshal : not thread-safe
func readAndUnmarshal(reader io.Reader) ([]byte, *relay_pb.Message, error) {
	sizeBuffer, err := readExact(reader, 8)
	if err != nil {
		return nil, nil, err
	}
	size := int(binary.LittleEndian.Uint64(sizeBuffer))
	dataBuffer, err := readExact(reader, size)
	if err != nil {
		return nil, nil, err
	}
	m := &relay_pb.Message{}
	err = proto.Unmarshal(dataBuffer, m)
	if err != nil {
		return nil, nil, err
	}
	return append(sizeBuffer, dataBuffer...), m, nil
}

func marshalAndWrite(writer io.Writer, m *relay_pb.Message) error {
	dataBuffer, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	sizeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuffer, uint64(len(dataBuffer)))
	n, err := writer.Write(append(sizeBuffer, dataBuffer...))
	if n != len(sizeBuffer)+len(dataBuffer) {
		return sizeError
	}
	return err
}
