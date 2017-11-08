package main

import (
	//	"gokit/log"
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
