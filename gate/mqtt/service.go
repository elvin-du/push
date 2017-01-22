package mqtt

import (
	"net"
	"time"
)

var (
	SERVICE_DEFAULT_OUT_CH_SIZE = 1024
)

type Service struct {
	Conn      net.Conn
	ClientId  string
	UserId    string
	Keepalive time.Duration

	outCh  chan []byte
	stopCh chan byte
}

func NewService(conn net.Conn) *Service {
	return &Service{
		Conn:   conn,
		outCh:  make(chan []byte, SERVICE_DEFAULT_OUT_CH_SIZE),
		stopCh: make(chan byte),
	}
}

func (s *Service) Run() {
	go s.ReadLoop()

	go s.WriteLoop()

	<-s.stopCh
}
