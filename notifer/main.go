package main

import (
	"push/notifer/service"
	"push/notifer/service/nsq/consumer"
)

func main() {
	service.Start()
	consumer.NewNsqConsumer(
		&SingleMsgHandler{},
		&BroadcastMsgHandler{},
	).Run()
}
