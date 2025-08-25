package relay

import (
	"log"
	"net"
	"sync"
)

type Hub interface {
	ListenAndServe() error
	Close() error
}

func NewHub(listenAddr string) (Hub, error) {
	listen, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &hub{
		listen:  listen,
		ln:      nil,
		connMap: sync.Map{},
	}, nil
}

type hub struct {
	listen  *net.TCPAddr
	ln      *net.TCPListener
	connMap sync.Map // map[name]*net.TCPAddr
}

func (h *hub) ListenAndServe() (err error) {
	h.ln, err = net.ListenTCP("tcp", h.listen)
	if err != nil {
		return err
	}

	log.Printf("listenning to %s\n", h.listen.String())
	for {
		conn, err := h.ln.AcceptTCP()
		if err != nil {
			return err
		}
		go h.handle(conn)
	}
}
func (h *hub) Close() error {
	return h.ln.Close()
}

func (h *hub) handle(conn *net.TCPConn) {
	var name string = ""
	defer func() {
		_ = conn.Close()
		if name != "" {
			h.connMap.Delete(name)
			log.Printf("peer [%s|%s] has been removed\n", name, conn.RemoteAddr().String())
		}
	}()
	for {
		buffer, m, err := readAndUnmarshal(conn)
		if err != nil {
			return
		}
		if name != m.Sender {
			if name != "" {
				h.connMap.Delete(name)
				log.Printf("peer [%s|%s] has been removed\n", name, conn.RemoteAddr().String())
			}
			name = m.Sender
			h.connMap.Store(name, conn)
			log.Printf("peer [%s|%s] has been registered\n", name, conn.RemoteAddr().String())
		}
		receiver := m.Receiver
		if val, loaded := h.connMap.Load(receiver); loaded {
			receiverConn := val.(*net.TCPConn)
			_, _ = receiverConn.Write(buffer)
		}
	}
}
