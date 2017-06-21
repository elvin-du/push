package main

import (
	"gokit/log"
	"push/notifer/service"
	"push/notifer/service/nsq/consumer"
)

func main() {
	log.Infoln("Notifer Runing!")
	service.Start()
	consumer.NewNsqConsumer(
		&SingleMsgHandler{},
		&BroadcastMsgHandler{},
	).Run()
}
