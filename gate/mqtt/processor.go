package mqtt

import (
	"gokit/log"
	"time"

	"github.com/surgemq/message"
)

func (s *Session) Process(msg message.Message) error {
	return s.readPacketCallback(s, msg)
	//	var err error = nil

	//	switch msg := msg.(type) {
	//	case *message.PubackMessage:
	//		return s.processPubAck(msg)
	//	case *message.DisconnectMessage:
	//		return s.processDisConn(msg)
	//	case *message.PingreqMessage:
	//		return s.processPingReq(msg)
	//	default:
	//	}

	//	return err
}

func (s *Session) processPubAck(msg *message.PubackMessage) error {
	//TODO pushlish成功，删除消息
	log.Debugf("got ack for %d,so remove it", msg.PacketId())
	log.Debugln(*msg)
	return nil
}

func (s *Session) processDisConn(msg *message.DisconnectMessage) error {
	//TODO 客户端要求断开链接，删除数据库
	log.Debugln(*msg)
	return nil
}

func (s *Session) processPingReq(msg *message.PingreqMessage) error {
	log.Debugln("ping came")
	//TODO 更新用户生命周期
	pingResp := message.NewPingrespMessage()
	err := s.WriteMsg(pingResp)
	if nil != err {
		log.Errorln(err)
		return err
	}

	s.SetTouchTime(time.Now().Unix())
	return nil
}
