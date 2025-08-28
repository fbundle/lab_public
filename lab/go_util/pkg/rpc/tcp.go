package rpc

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	DEFAULT_TCP_TIMEOUT = 10 * time.Second
)

type TCPServer interface {
	ListenAndServe(ctx context.Context, dispatcher Dispatcher, msgIO MessageIO) error
	Close() error
}

func TCPTransport(ctx context.Context, addr string, msgIO MessageIO) TransportFunc {
	return func(b []byte) ([]byte, error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer conn.Close()

		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(DEFAULT_TCP_TIMEOUT)
		}
		err = conn.SetDeadline(deadline)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		err = msgIO.Write(ctx, conn, b)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		b, err = msgIO.Read(ctx, conn)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return b, nil
	}
}

type tcpServer struct {
	listener net.Listener
}

func NewTCPServer(bindAddr string) (TCPServer, error) {
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return nil, err
	}
	return &tcpServer{
		listener: listener,
	}, nil
}

func (s *tcpServer) Close() error {
	return s.listener.Close()
}

func (s *tcpServer) ListenAndServe(ctx context.Context, dispatcher Dispatcher, msgIO MessageIO) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(ctx, dispatcher, msgIO, conn)
	}
}
func (s *tcpServer) handleConn(ctx context.Context, dispatcher Dispatcher, msgIO MessageIO, conn net.Conn) {
	defer conn.Close()
	err := conn.SetDeadline(time.Now().Add(DEFAULT_TCP_TIMEOUT))
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := msgIO.Read(ctx, conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err = dispatcher.Handle(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = msgIO.Write(ctx, conn, b)
	if err != nil {
		fmt.Println(err)
		return
	}
}
