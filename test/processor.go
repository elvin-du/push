package main

import (
	"encoding/json"
	"log"

	"github.com/surgemq/message"
)

type PublishMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type PubAck struct {
	Kind    int         `json:"kind"`
	Content interface{} `json:"content"`
}

type PublishAck struct {
	MsgID string `json:"msg_id"`
}

func Process(msg message.Message) error {
	var err error = nil

	switch msg := msg.(type) {
	case *message.ConnackMessage:
		return processConnAck(msg)
	case *message.PublishMessage:
		return processPub(msg)
	case *message.PingrespMessage:
		return processPingResp(msg)
	default:
	}

	return err
}

func processConnAck(msg *message.ConnackMessage) error {
	log.Println("processConnAck", *msg)
	return nil
}

func processPub(msg *message.PublishMessage) error {
	log.Println("processPub", string(msg.Payload()))
	pubMsg := PublishMessage{}
	err := json.Unmarshal(msg.Payload(), &pubMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	ret := PubAck{}
	ack := PublishAck{}
	ack.MsgID = pubMsg.ID
	ret.Content = ack
	ret.Kind = 1

	bin, err := json.Marshal(ret)
	if nil != err {
		log.Println(err)
		return err
	}

	ackMsg := message.NewPublishMessage()
	ackMsg.SetQoS(message.QosAtLeastOnce)
	ackMsg.SetPacketId(1)
	ackMsg.SetTopic([]byte("*"))
	ackMsg.SetPayload(bin)

	err = Send(ackMsg)
	if nil != err {
		log.Println(err)
		return err
	}
	return nil
}

func processPingResp(msg *message.PingrespMessage) error {
	log.Println("processPingResp", *msg)
	return nil
}
