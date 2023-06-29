package gotcp

import (
	"net"
)

type Server interface {
	Run() error
	Stop()
}

func NewServer(handler EventHandler, addr string) Server {
	return &server{
		listener: nil,
		addr:     addr,
		handler:  handler,
	}
}

type server struct {
	listener net.Listener
	addr     string
	handler  EventHandler
}

func (s *server) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		session := NewSession(conn, s.handler)
		go session.Start()
	}
}

func (s *server) Stop() {
	s.listener.Close()
}
