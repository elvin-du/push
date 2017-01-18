package mqtt

import (
	"errors"
	"log"
)

var (
	E_READ_ERROR = errors.New("read error")
)

func (s *Service) WriteLoop() error {
	for {
		select {
		case out := <-s.outCh:
			n, err := s.Conn.Write(out)
			if nil != err {
				log.Println(err)
				return err
			}
			log.Println("wrote number:", n)
		}
	}
	return E_WRITE_ERROR
}

func (s *Service) Push(content []byte) {
	s.outCh <- content
}
