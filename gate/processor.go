package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	"push/common/model"
	"push/gate/mqtt"
	"push/gate/service/session"
	"time"

	gateMsg "push/gate/message"

	"github.com/surgemq/message"
)

type PublishAck struct {
	MsgID string `json:"msg_id"`
}

func Dispatch(ses *mqtt.Session, msg message.Message) error {
	switch msg := msg.(type) {
	case *message.DisconnectMessage:
		return processDisConn(ses, msg)
	case *message.PingreqMessage:
		return processPingReq(ses, msg)
	case *message.PublishMessage:
		return processPublish(ses, msg)
	case *message.ConnectMessage:
		return processConnect(ses, msg)
	}

	log.Errorf("unsupport msg type:%d,%s", msg.Type(), msg.Name())
	err := errors.New("unsupport msg type")
	return err
}

func processPublish(ses *mqtt.Session, msg *message.PublishMessage) error {
	log.Debugf("processPublish:%+v", *msg)

	ack := &PublishAck{}
	err := json.Unmarshal(msg.Payload(), ack)
	if nil != err {
		log.Errorln(err)
		return err
	}
	log.Debugf("got ack for %+v", ack)

	if gateMsg.DefaultMessageManager.IsExist(ack.MsgID) {
		log.Infof("remove msg:%s from messageManager", ack.MsgID)
		gateMsg.DefaultMessageManager.Delete(ack.MsgID)
	} else {
		log.Infof("remove msg:%s from DB", ack.MsgID)
		err = model.OfflineMsgModel().Delete(ack.MsgID)
		if nil != err {
			log.Errorln(err, "msg_id:", ack.MsgID)
			return err
		}
	}

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
	users := defaultServer.GetByID(ses.ID)
	for _, u := range users {
		session.Touch(u.AppID, u.RegID)
	}

	return nil
}

func processConnect(ses *mqtt.Session, msg *message.ConnectMessage) (err error) {
	log.Debugln("connect came")
	connAckMsg := message.NewConnackMessage()

	defer func() {
		//回应connect消息
		if nil != err {
			connAckMsg.SetReturnCode(message.ErrNotAuthorized)
			err2 := ses.SendMsg(connAckMsg)
			if nil != err {
				log.Error(err2)
			}

			//TODO 直接关闭连接?
			ses.Close(err)
		}
	}()

	//合法性检验
	var u *User
	u, err = Auth(ses, msg)
	if nil != err {
		log.Error(err)
		return err
	}
	u.SetTouchTime(time.Now().Unix())

	//连接成功
	connAckMsg.SetReturnCode(message.ConnectionAccepted)
	err = ses.SendMsg(connAckMsg)
	if nil != err {
		log.Error(err)
		return err
	}

	//发送离线消息
	u.CheckOfflineMsgs()
	return nil
}
