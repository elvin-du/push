package mqtt

import (
	"net"
	"time"
)

var (
	SERVICE_DEFAULT_OUT_CH_SIZE = 1024
)

type Service struct {
	Conn         net.Conn
	ClientId     string
	UserId       string
	readTimeout  time.Duration
	writeTimeout time.Duration
	touchTime    int64

	outCh chan []byte
}

func NewService(conn net.Conn) *Service {
	return &Service{
		Conn:  conn,
		outCh: make(chan []byte, SERVICE_DEFAULT_OUT_CH_SIZE),
	}
}

func (s *Service) Run() {
	go s.ReadLoop()
	go s.WriteLoop()

}

func (s *Service) SetReadTimeout(d time.Duration) {
	s.readTimeout = d
}

func (s *Service) SetWriteTimeout(d time.Duration) {
	s.writeTimeout = d
}

func (s *Service) SetTouchTime(t int64) {
	s.touchTime = t
}

func (s *Service) IsAlive(t int64) bool {
	return time.Now().Unix()-s.touchTime < t
}
