package main

import (
	"encoding/json"
	"gokit/log"
	"net/http"
	"push/rest_api/service/nsq"

	"gopkg.in/mgo.v2/bson"
)

type push struct{}

var _push *push

func (*push) Push(ctx *Context) {
	bin, err := ctx.GetRawData()
	if nil != err {
		ctx.AbortWithError(400, err)
		return
	}

	log.Debugf("Msg:%+v", string(bin))
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
		msg.ID = bson.NewObjectId().Hex()
		info := &Info{Message: msg, AppID: ctx.AppID}
		bin, err = json.Marshal(info)
		if nil != err {
			ctx.AbortWithError(400, err)
			return
		}

		err = SaveMsg(info)
		if nil != err {
			ctx.AbortWithError(500, err)
			return
		}

		go nsq.SingleProducer.Publish(bin)
	default:
		ctx.AbortWithError(400, REQUEST_DATA_INVALID)
		return
	}

	ctx.Status(http.StatusOK)
}
