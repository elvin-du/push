package main

import (
	"errors"
	"log"

	"github.com/nsqio/go-nsq"
)

type BroadcastMsgHandler struct{}

func (b *BroadcastMsgHandler) Process(i interface{}) error {
	msg, ok := i.(*nsq.Message)
	if !ok {
		err := errors.New("i is not *nsq.Message")
		log.Println(err)
		return err
	}
	log.Println(string(msg.Body))

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

	return nil
}
