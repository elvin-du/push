package main

import (
	"gokit/log"
	"push/notifer/service/nsq/consumer"
)

func main() {
	log.Infoln("Notifer Runing!")
	consumer.NewNsqConsumer(
		&SingleMsgHandler{},
		&BroadcastMsgHandler{},
	).Run()
}
