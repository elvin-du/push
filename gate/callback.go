package main

import (
	"gokit/log"
	"io"
	gateMsg "push/gate/message"
	"push/gate/mqtt"

	"github.com/surgemq/message"
)

func OnSend(ses *mqtt.Session, data []byte) error {
	//	msg := message.NewPublishMessage()
	//	_, err := msg.Decode(data)
	//	if nil != err {
	//		log.Errorln(err)
	//		return err
	//	}

	return nil
}

func OnRead(ses *mqtt.Session, msg message.Message) error {
	return Dispatch(ses, msg)
}

func OnClose(ses *mqtt.Session, err error) {
	u := defaultServer.RemoveByID(ses.ID)
	var appID, regID string
	if nil != u {
		appID = u.AppID
		regID = u.RegID
	}
	log.Infof("remove user(app_id:%s,reg_id:%s) session", appID, regID)

	if io.EOF == err {
		log.Infof("app_id:%s,reg_id:%s session close,err:%s", appID, regID, err.Error())
	} else {
		log.Errorf("app_id:%s,reg_id:%s session close,err:%s", appID, regID, err.Error())
	}
	//TODO
	gateMsg.DefaultMessageManager.SyncByAccount(ses.ID)
}
