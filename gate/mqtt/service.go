package mqtt

import (
	"net"
	"sync"
)

var (
	SERVICE_DEFAULT_OUT_CH_SIZE = 1024
)

type Service struct {
	Conn      net.Conn
	wgStarted sync.WaitGroup
	ClientId  string
	UserId    string

	outCh chan []byte
}

func NewService(conn net.Conn) *Service {
	return &Service{
		Conn:  conn,
		outCh: make(chan []byte, SERVICE_DEFAULT_OUT_CH_SIZE),
	}
}

func (s *Service) Start() error {
	go s.ReadLoop()

	go s.WriteLoop()

	s.wgStarted.Wait()
	return nil
}
