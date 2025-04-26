package relay

import (
	"errors"
	"fmt"
	"github.com/khanh-nguyen-code/go_util/pkg/relay/proto/gen/relay_pb"
	"net"
)

var (
	ErrNotInit = errors.New("not_init")
)

type Peer interface {
	Write(m *relay_pb.Message) error
	DialAndServe(serve func(m *relay_pb.Message)) error
	Close() error
}

func NewPeer(name string, listenAddr string, relayAddr string) (Peer, error) {
	listen, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	relay, err := net.ResolveTCPAddr("tcp", relayAddr)
	if err != nil {
		return nil, err
	}
	return &peer{
		name:   name,
		listen: listen,
		relay:  relay,
		conn:   nil,
	}, nil
}

type peer struct {
	name   string
	listen *net.TCPAddr
	relay  *net.TCPAddr
	conn   *net.TCPConn
}

func (p *peer) Write(m *relay_pb.Message) error {
	m.Sender = p.name
	if p.conn == nil {
		err := ErrNotInit
		fmt.Printf("[peer_%s] write error: %v\n", p.name, err)
		return err
	}
	err := marshalAndWrite(p.conn, m)
	if err != nil {
		fmt.Printf("[peer_%s] write error: %v\n", p.name, err)
		return err
	}
	return err
}

func (p *peer) DialAndServe(serve func(*relay_pb.Message)) (err error) {
	fmt.Printf("[peer_%s] dialing %s\n", p.name, p.relay.String())
	p.conn, err = net.DialTCP("tcp", p.listen, p.relay)
	if err != nil {
		fmt.Printf("[peer_%s] dial error: %v\n", p.name, err)
		return err
	}
	defer p.conn.Close()
	fmt.Printf("[peer_%s] connected to relay %s\n", p.name, p.relay.String())
	err = p.Write(&relay_pb.Message{})
	if err != nil {
		return err
	}
	for {
		_, m, err := readAndUnmarshal(p.conn)
		if err != nil {
			fmt.Printf("[peer_%s] read error: %v\n", p.name, err)
			return err
		}
		serve(m)
	}
}

func (p *peer) Close() error {
	if p.conn == nil {
		err := ErrNotInit
		fmt.Printf("[agent_%s] close error: %v\n", p.name, err)
		return err
	}
	return p.conn.Close()
}
