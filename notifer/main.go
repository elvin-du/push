package main

import (
	"hscore/log"
	"push/notifer/nsq/consumer"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	consumer.NewNsqConsumer(
		&SingleMsgHandler{},
		&BroadcastMsgHandler{},
	).Run()
}
