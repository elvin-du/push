package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	gateCli "push/gate/client"
	"push/notifer/service/db"
	"push/pb"

	"github.com/nsqio/go-nsq"
)

var (
	GATE_RPC_PORT = ":50002"
)

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
		return errors.New("not found session")
	}

	_, err = gateCli.Push(
		ses.GateServerIP,
		ses.GateServerPort,
		&pb.GatePushRequest{
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
