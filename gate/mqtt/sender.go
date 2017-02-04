package mqtt

import (
	"errors"
	"hscore/log"
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
			err := s.SetWriteDeadline(s.Keepalive)
			if nil != err {
				log.Error(err)
				continue
			}

			n, err := s.Conn.Write(out)
			if nil != err {
				log.Error(err)
				return err
			}
			log.Debugln("wrote number:", n)
		}
	}
	return E_WRITE_ERROR
}

func (s *Service) Push(content []byte) error {
	pubMsg := message.NewPublishMessage()
	pubMsg.SetQoS(message.QosAtLeastOnce)
	pubMsg.SetPacketId(uint16(rand.Uint32()))
	pubMsg.SetTopic([]byte("*"))
	pubMsg.SetPayload(content)

	return s.Write(pubMsg)
}

func (s *Service) Write(msg message.Message) error {
	log.Debugln("write:", msg.Desc())
	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err {
		log.Error(err)
		return err
	}

	s.outCh <- buf[:n]

	return nil
}
