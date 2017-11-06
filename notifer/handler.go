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
	msg, ok := i.(*nsq.Message)
	if !ok {
		err := errors.New("parameter is not *nsq.Message type")
		log.Errorln(err)
		return err
	}
	log.Debugln(string(msg.Body))

	data := Message{}
	err := json.Unmarshal(msg.Body, &data)
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
		log.Errorf("not found session by key :%s", data.Key())
		offlineMsg := &model.OfflineMsg{}
		offlineMsg.AppID = data.AppID
		offlineMsg.ClientID = data.ClientID
		offlineMsg.Content = data.Content
		offlineMsg.Extra = data.Extra
		offlineMsg.Kind = data.Kind
		offlineMsg.ID = data.ID

		err = model.OfflineMsgModel().Insert(offlineMsg)
		if nil != err {
			log.Errorln(err)
			return err
		}

		return errors.New("not found session")
	}

	_, err = gateCli.Push(
		ses.GateServerIP,
		ses.GateServerPort,
		&pb.GatePushRequest{
			ID:       data.ID,
			AppID:    data.AppID,
			ClientId: data.ClientID,
			Content:  data.Content,
			Kind:     data.Kind,
			Extra:    data.Extra,
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
