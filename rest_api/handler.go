package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"push/rest_api/service/nsq"
)

type push struct{}

var _push *push

func (*push) Push(ctx *Context) {
	bin, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		ctx.AbortWithError(400, err)
		return
	}

	msg, err := ValidMessage(bin)
	if nil != err {
		ctx.AbortWithError(400, err)
		return
	}

	switch msg.Platform {
	case PUSH_PLATFORM_IOS:
		err = IOSPush(ctx.AppID, msg.ClientID, msg.Content, msg.IsProduction)
		if nil != err {
			ctx.AbortWithError(500, err)
			return
		}
	case PUSH_PLATFORM_ANDROID:
		info := &Info{Message: msg, AppID: ctx.AppID}
		bin, err = json.Marshal(info)
		if nil != err {
			ctx.AbortWithError(400, err)
			return
		}

		go nsq.SingleProducer.Publish(bin)
	default:
		ctx.AbortWithError(400, REQUEST_DATA_INVALID)
		return
	}

	ctx.Status(http.StatusOK)
}
