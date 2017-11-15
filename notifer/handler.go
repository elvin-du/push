package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	"push/common/db"
	"push/common/model"
	gateCli "push/gate/client"
	"push/pb"

	"github.com/nsqio/go-nsq"
)

type SingleMsgHandler struct{}

func (b *SingleMsgHandler) Process(i interface{}) error {
	nsqMsg, ok := i.(*nsq.Message)
	if !ok {
		err := errors.New("parameter is not *nsq.Message type")
		log.Errorln(err)
		return err
	}
	log.Debugln(string(nsqMsg.Body))

	data := Message{}
	err := json.Unmarshal(nsqMsg.Body, &data)
	if nil != err {
		log.Errorln(err)
		return err
	}

	var ses session
	err = db.MainRedis().HGETALL(data.Key(), &ses)
	if nil != err {
		log.Errorln(err)
		return err
	}
	log.Debugln("session:", ses)
	if ses.GateServerIP == "" && "" == ses.GateServerPort {
		log.Errorf("not found session by key %s", data.Key())
		msg := &model.Message{}
		msg.AppID = data.AppID
		msg.RegID = data.RegID
		msg.Content = data.Content
		msg.Extra = data.Extra
		msg.Kind = data.Kind
		msg.ID = data.ID

		err = model.MessageModel().Insert(msg)
		if nil != err {
			log.Errorln(err)
			return err
		}

		//不在线不需要nsq消息重发
		return nil
	}

	_, err = gateCli.Push(
		ses.GateServerIP,
		ses.GateServerPort,
		&pb.GatePushRequest{
			ID:      data.ID,
			AppID:   data.AppID,
			RegID:   data.RegID,
			Content: data.Content,
			Kind:    data.Kind,
			Extra:   data.Extra,
		})
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

type BroadcastMsgHandler struct{}

//TODO do it later
func (b *BroadcastMsgHandler) Process(i interface{}) error {
	//	msg, ok := i.(*nsq.Message)
	//	if !ok {
	//		err := errors.New("i is not *nsq.Message")
	//		log.Println(err)
	//		return err
	//	}
	//	log.Println(string(msg.Body))

	//	data := Message{}
	//	err := json.Unmarshal(msg.Body, &data)
	//	if nil != err {
	//		log.Println(err)
	//		return err
	//	}

	return nil
}
