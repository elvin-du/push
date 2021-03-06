package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	"push/common/model"
	"push/gate/mqtt"
	"push/gate/service/session"
	"time"

	"github.com/surgemq/message"
)

const (
	CLIENT_MESSAGE_KIND_PUB_ACK = 1
	CLIENT_MESSAGE_KIND_SIGNOUT = 2
)

type PublishAck struct {
	MsgID string `json:"msg_id"`
}

type SignOut struct {
	AppID string `json:"app_id"`
	RegID string `json:"reg_id"`
}

type ClientMessage struct {
	Kind    int         `json:"kind"` //1:publish ack 2:signout
	Content interface{} `json:"content"`
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

	log.Errorf("unsupport message type:%d,%s", msg.Type(), msg.Name())
	err := errors.New("unsupport message type")
	return err
}

func processPublish(ses *mqtt.Session, msg *message.PublishMessage) error {
	log.Debugf("process publish msg:%+v", string(msg.Payload()))

	cliMsg := &ClientMessage{}
	err := json.Unmarshal(msg.Payload(), cliMsg)
	if nil != err {
		log.Errorln(err)
		return err
	}

	switch cliMsg.Kind {
	case CLIENT_MESSAGE_KIND_PUB_ACK:
		return handlePubAck(cliMsg.Content)
	case CLIENT_MESSAGE_KIND_SIGNOUT:
		return handleSignOut(cliMsg.Content)
	}

	log.Errorln("invalid message kind")
	return errors.New("Invalid message kind")
}

func handleSignOut(data interface{}) error {
	bin, err := json.Marshal(data)
	if nil != err {
		log.Errorln(err)
		return err
	}
	log.Debugf("signout message data:%+v", string(bin))
	out := &SignOut{}
	err = json.Unmarshal(bin, out)
	if nil != err {
		log.Errorln(err)
		return err
	}

	log.Infof("remove user(app_id:%s,reg_id:%s) cache from mem", out.AppID, out.RegID)
	defaultServer.RemoveUser(out.AppID, out.RegID)
	return nil
}

func handlePubAck(data interface{}) error {
	bin, err := json.Marshal(data)
	if nil != err {
		log.Errorln(err)
		return err
	}

	ack := &PublishAck{}
	err = json.Unmarshal(bin, ack)
	if nil != err {
		log.Errorln(err)
		return err
	}
	log.Debugf("publish message ack data: %+v", ack)

	//	if gateMsg.DefaultMessageManager.IsExist(ack.MsgID) {
	//		log.Infof("remove msg:%s from messageManager", ack.MsgID)
	//		gateMsg.DefaultMessageManager.Delete(ack.MsgID)
	//	} else {
	log.Infof("remove message msg_id:%s from DB", ack.MsgID)
	err = model.MessageModel().Delete(ack.MsgID)
	if nil != err {
		log.Errorf("delete message(id:%s) failed,err:%s ", ack.MsgID, err.Error())
		return err
	}
	//	}

	return nil
}

func processDisConn(ses *mqtt.Session, msg *message.DisconnectMessage) error {
	//TODO 客户端要求断开链接，删除数据库
	log.Debugln(*msg)
	return nil
}

func processPingReq(ses *mqtt.Session, msg *message.PingreqMessage) error {
	log.Infof("ping message came,session_id:%s", ses.ID)

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
	log.Debugln("connect message came")
	connAckMsg := message.NewConnackMessage()

	defer func() {
		//回应connect消息
		if nil != err {
			connAckMsg.SetReturnCode(message.ErrNotAuthorized)
			err2 := ses.SendMsg(connAckMsg)
			if nil != err2 {
				log.Error(err2)
			}

			//TODO 直接关闭连接?
			//ses.Close(err)
		}
	}()

	//合法性检验
	var u *User
	u, err = Auth(ses, msg)
	if nil != err {
		log.Error(err)
		return err
	}

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
