package main

import (
	"fmt"
	"log"

	"github.com/surgemq/message"
	//	"github.com/surgemq/surgemq/service"
)

func Send(msg message.Message) error {
	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err && n == len(buf) {
		log.Println(err)
		return err
	}

	n, err = conn.Write(buf)
	if nil != err {
		log.Println(err)
		return err
	}

	if n != len(buf) {
		err = fmt.Errorf("expected len:%d,got:%d", len(buf), n)
		log.Println(err)
		return err
	}

	return nil
}
