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
	dst := make([]byte, 1)

	pubMsg := &message.PublishMessage{}
	pubMsg.SetQoS(message.QosAtLeastOnce)
	pubMsg.SetPacketId(uint16(rand.Uint32()))
	pubMsg.SetPayload(content)
	pubMsg.SetRemainingLength(int32(len(content)))
	n, err := pubMsg.Encode(dst)
	if nil != err {
		log.Println(err)
		return
	}

	dst = dst[:n]

	s.outCh <- dst
}
