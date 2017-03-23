package main

import (
	"hscore/log"

	"github.com/surgemq/message"
)

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
	log.Println(*msg)
	return nil
}

func processPub(msg *message.PublishMessage) error {
	log.Println(*msg)
	return nil
}

func processPingResp(msg *message.PingrespMessage) error {
	log.Println(*msg)
	return nil
}
