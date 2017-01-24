package mqtt

import (
	"log"

	"github.com/surgemq/message"
)

func (s *Service) Process(msg message.Message) error {
	var err error = nil

	switch msg := msg.(type) {
	case *message.PubackMessage:
		return s.processPubAck(msg)
	case *message.DisconnectMessage:
		return s.processDisConn(msg)
	case *message.PingreqMessage:
		return s.processPingReq(msg)
	default:
	}

	return err
}

func (s *Service) processPubAck(msg *message.PubackMessage) error {
	//TODO pushlish成功，删除消息
	log.Println(*msg)
	return nil
}

func (s *Service) processDisConn(msg *message.DisconnectMessage) error {
	//TODO 客户端要求断开链接，删除数据库
	log.Println(*msg)
	return nil
}

func (s *Service) processPingReq(msg *message.PingreqMessage) error {
	log.Println("ping came")
	pingResp := message.NewPingrespMessage()
	s.Write(pingResp)
	return nil
}
