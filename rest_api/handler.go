package main

import (
	"log"
	"net/http"
	"push/data/client"
	"push/meta"
	"push/rest_api/nsq/producer"

	"github.com/gin-gonic/gin"
)

func init() {
	sigleProducer.Start()
	broadcastProducer.Start()
}

var (
	NSQD_ADDR           = "127.0.0.1:4150"
	NSQ_TOPIC_SINGLE    = "topic_push_single"
	NSQ_TOPIC_BROADCAST = "topic_push_broadcast"
)

var (
	sigleProducer     = producer.NewNsqProducer(NSQD_ADDR, NSQ_TOPIC_SINGLE)
	broadcastProducer = producer.NewNsqProducer(NSQD_ADDR, NSQ_TOPIC_BROADCAST)
)

type push struct{}

var _push *push

func (*push) Push(ctx *gin.Context) {
	resp, err := client.Online(&meta.DataOnlineRequest{UserId: "123", IP: "12.12.151.2"})
	if nil != err {
		log.Println(err)
	}
	log.Println(resp.String())
	go sigleProducer.Publish([]byte(resp.String()))
	go broadcastProducer.Publish([]byte(resp.String()))

	ctx.String(http.StatusOK, "hi push")
}
