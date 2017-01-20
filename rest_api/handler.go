package main

import (
	"errors"
	"net/http"
	"push/rest_api/nsq/producer"

	"io/ioutil"

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
	bin, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		ctx.AbortWithError(400, errors.New("parameter invalid"))
		return
	}

	go sigleProducer.Publish(bin)
	go broadcastProducer.Publish(bin)

	ctx.Status(http.StatusOK)
}
