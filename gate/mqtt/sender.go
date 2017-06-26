package mqtt

import (
	"errors"
	"gokit/log"
	"time"

	"github.com/surgemq/message"
)

var (
	E_READ_ERROR = errors.New("read error")
)

func (s *Service) WriteLoop() error {
	for {
		select {
		case out := <-s.outCh:
			err := s.Send(out)
			if nil != err {
				log.Errorln(err)
				//TODO 直接断开连接？等待心跳机制检查时删除本服务？
				return err
			}
		}
	}

	return E_WRITE_ERROR
}

func (s *Service) Push(packetId uint16, content []byte) error {
	pubMsg := message.NewPublishMessage()
	pubMsg.SetQoS(message.QosAtLeastOnce)
	pubMsg.SetPacketId(packetId)
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
	log.Debugln("n:", n, string(buf))

	s.outCh <- buf[:n]

	return nil
}

func (s *Service) SendMsg(msg message.Message) error {
	log.Debugln("write:", msg.Desc())
	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err {
		log.Error(err)
		return err
	}

	return s.Send(buf[:n])
}

func (s *Service) Send(data []byte) error {
	if s.writeTimeout > 0 {
		s.Conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	} else {
		s.Conn.SetWriteDeadline(time.Time{})
	}

	_, err := s.Conn.Write(data)
	if nil != err {
		log.Error(err)
		return err
	}

	return nil
}
