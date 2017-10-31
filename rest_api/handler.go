package main

import (
	"io/ioutil"
	"net/http"
	"push/rest_api/service/nsq"

	"github.com/gin-gonic/gin"
)

type push struct{}

var _push *push

func (*push) Push(ctx *gin.Context) {
	err := Auth(ctx)
	if nil != err {
		ctx.AbortWithError(400, err)
		return
	}

	bin, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		ctx.AbortWithError(400, err)
		return
	}

	err = ValidMessage(bin)
	if nil != err {
		ctx.AbortWithError(400, err)
		return
	}

	go nsq.SingleProducer.Publish(bin)

	ctx.Status(http.StatusOK)
}
