package main

import (
	"gokit/log"
	"io"
	//	gateMsg "push/gate/message"
	"push/gate/mqtt"

	"github.com/surgemq/message"
)

func OnSend(ses *mqtt.Session, data []byte) error {
	msg, err := mqtt.NewMessage(data)
	if nil != err {
		log.Errorln(err)
		return err
	}

	switch msg := msg.(type) {
	case *message.PublishMessage:
		log.Infof("push sucess,data:%+v", *msg)
	default:
		return nil //TODO
	}

	return nil
}

func OnRead(ses *mqtt.Session, msg message.Message) error {
	return Dispatch(ses, msg)
}

func OnClose(ses *mqtt.Session, err error) {
	users := defaultServer.RemoveByID(ses.ID)
	for _, u := range users {
		if nil != u {
			log.Infof("remove user(app_id:%s,reg_id:%s,session_id:%s) session", u.AppID, u.RegID, u.ID)
			if io.EOF == err {
				log.Infof("app_id:%s,reg_id:%s,session_id:%s session close,err:%s", u.AppID, u.RegID, u.ID, err.Error())
			} else {
				log.Errorf("app_id:%s,reg_id:%s,session_id:%s session close,err:%s", u.AppID, u.RegID, u.ID, err.Error())
			}
			//			gateMsg.DefaultMessageManager.SyncByAccount(u.Key())
		} else {
			log.Errorf("session has nil user,ID:%s", ses.ID)
		}
	}
}
