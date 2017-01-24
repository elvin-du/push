package main

import (
	"encoding/json"
	"errors"
	"log"
	dataCli "push/data/client"
	gateCli "push/gate/client"
	"push/meta"

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
		err := errors.New("i is not *nsq.Message")
		log.Println(err)
		return err
	}
	log.Println(string(msg.Body))

	data := Message{}
	err := json.Unmarshal(msg.Body, &data)
	if nil != err {
		log.Println(err)
		return err
	}

	resp, err := dataCli.GetClientInfo(&meta.GetClientInfoRequest{UserId: data.UserId})
	if nil != err {
		log.Println(err)
		return err
	}
	log.Printf("gate info:%+v", resp)

	for _, item := range resp.GetItems() {
		if item.Platform == "android" {
			log.Println(item.IP)
			_, err = gateCli.Push(item.IP, GATE_RPC_PORT, &meta.GatePushRequest{
				UserId:  data.UserId,
				Content: data.Content,
				Kind:    data.Kind,
				Extra:   data.Extra,
			})
			if nil != err {
				log.Println(err)
				return err
			}
		} else if "ios" == item.Platform {
			//TODO
		}
	}

	return nil
}
