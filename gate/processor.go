package main

import (
	"errors"
	"gokit/log"
	"push/gate/mqtt"
	"push/gate/service/session"
	"time"

	"github.com/surgemq/message"
)

func Dispatch(ses *mqtt.Session, msg message.Message) error {
	switch msg := msg.(type) {
	case *message.PubackMessage:
		return processPubAck(ses, msg)
	case *message.DisconnectMessage:
		return processDisConn(ses, msg)
	case *message.PingreqMessage:
		return processPingReq(ses, msg)
	}

	log.Errorf("unsupport msg type:%d,%s", msg.Type(), msg.Name())
	err := errors.New("unsupport msg type")
	return err
}

func processPubAck(ses *mqtt.Session, msg *message.PubackMessage) error {
	//TODO pushlish成功，删除消息
	log.Debugf("got ack for %d,so remove it", msg.PacketId())
	log.Debugln(*msg)
	return nil
}

func processDisConn(ses *mqtt.Session, msg *message.DisconnectMessage) error {
	//TODO 客户端要求断开链接，删除数据库
	log.Debugln(*msg)
	return nil
}

func processPingReq(ses *mqtt.Session, msg *message.PingreqMessage) error {
	log.Debugln("ping came")

	pingResp := message.NewPingrespMessage()
	err := ses.WriteMsg(pingResp)
	if nil != err {
		log.Errorln(err)
		return err
	}

	//更新用户生命周期
	ses.SetTouchTime(time.Now().Unix())
	session.Touch(ses.AppID, ses.ClientID)
	return nil
}
