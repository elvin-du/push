package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	"push/common/db"
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
	log.Debugf("received data:%s", string(nsqMsg.Body))

	data := NewMessage()
	err := json.Unmarshal(nsqMsg.Body, data)
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
	log.Debugf("msg_id:%s,session:%+v", data.ID, ses)
	if ses.GateServerIP == "" && "" == ses.GateServerPort {
		log.Errorf("session not found by key:%s,msg_id:%s", data.Key(), data.ID)
		//		msg := &model.Message{}
		//		msg.AppID = data.AppID
		//		msg.Content = data.Content
		//		msg.Extras = data.Extras
		//		msg.ID = data.ID
		//		msg.RegID = data.RegID
		//		msg.TTL = data.TTL
		//		msg.CreatedAt = util.Timestamp()
		//		err = model.MessageModel().Insert(msg)
		//		if nil != err {
		//			log.Errorln(err)
		//			return err
		//		}

		//不在线不需要nsq消息重发
		return nil
	}

	bin, err := json.Marshal(data.Extras)
	if nil != err {
		log.Errorln(err)
		return err
	}

	_, err = gateCli.Push(
		ses.GateServerIP,
		ses.GateServerPort,
		&pb.GatePushRequest{
			ID:      data.ID,
			AppID:   data.AppID,
			RegID:   data.RegID,
			Content: data.Content,
			TTL:     data.TTL,
			Extras:  string(bin),
		})
	if nil != err {
		log.Errorf("push failed,msg_id:%s,err:%s", data.ID, err.Error())
		//不需要nsq消息重发
		return nil
	}

	log.Infof("push success,msg_id:%s", data.ID)
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
