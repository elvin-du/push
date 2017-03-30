package main

import (
	"hscore/log"
	"push/notifer/service/nsq/consumer"
)

func main() {
	log.Infoln("Notifer Runing!")
	consumer.NewNsqConsumer(
		&SingleMsgHandler{},
		&BroadcastMsgHandler{},
	).Run()
}
