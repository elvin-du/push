package main

import (
	"log"
	"time"

	"github.com/surgemq/message"
)

func Ping() {
	for {
		pingMsg := message.NewPingreqMessage()
		err := Send(pingMsg)
		if nil != err {
			log.Println(err)
			return
		}
		time.Sleep(time.Second * 60)
	}
}
