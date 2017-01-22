package mqtt

import (
	"errors"
	"log"
	"math/rand"

	"github.com/surgemq/message"
)

var (
	E_READ_ERROR = errors.New("read error")
)

func (s *Service) WriteLoop() error {
	for {
		select {
		case out := <-s.outCh:
			log.Println("should send len", len(out))
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

func (s *Service) Push(content []byte) error {
	pubMsg := &message.PublishMessage{}
	pubMsg.SetQoS(message.QosAtLeastOnce)
	pubMsg.SetPacketId(uint16(rand.Uint32()))
	pubMsg.SetPayload(content)

	return s.Write(pubMsg)
}

func (s *Service) Write(msg message.Message) error {
	log.Println("write:", msg.Desc())
	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err {
		log.Println(err)
		return err
	}

	s.outCh <- buf[:n]

	return nil
}
