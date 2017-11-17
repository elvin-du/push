package main

import (
	"encoding/json"
	"gokit/log"
	"net/http"
	"push/common/model"
	"push/rest_api/client"
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
	notification, err := ValidateNotification(bin)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	resp := make(map[string]interface{})
	if client.C_AUDIENCE_ALL == notification.Audience {
		//TODO
	} else {
		reg, err := model.RegistryModel().Get(notification.Audience)
		if nil != err {
			log.Errorln(err)
			ctx.AbortWithError(400, err)
			return
		}

		switch reg.Platform {
		case client.C_PLATFORM_IOS:
			if nil != notification.IosNotification {
				for k, v := range notification.IosNotification.Extras {
					notification.AddExtra(k, v)
				}
			}

			apnsID, err := IOSPush(ctx.AppID, reg.DevToken, notification.Alert, notification.IosNotification.Sound,
				notification.IosNotification.Production, notification.IosNotification.Badge, notification.TTL, notification.Extras)
			if nil != err {
				ctx.AbortWithError(500, err)
				return
			}
			resp["msg_id"] = apnsID
		case client.C_PLATFORM_ANDRID:
			msg := NewMessage()
			msg.ID = bson.NewObjectId().Hex()
			msg.Content = notification.Alert
			if nil != notification.AndroidNotification {
				for k, v := range notification.AndroidNotification.Extras {
					notification.AddExtra(k, v)
				}
			}

			msg.Extras = notification.Extras
			msg.AppID = ctx.AppID
			msg.TTL = notification.TTL
			msg.RegID = notification.Audience
			bin, err = json.Marshal(msg)
			if nil != err {
				log.Errorln(err)
				ctx.AbortWithError(400, err)
				return
			}

			err = SaveMsg(msg)
			if nil != err {
				log.Errorln(err)
				ctx.AbortWithError(500, err)
				return
			}
			resp["msg_id"] = msg.ID
			go nsq.SingleProducer.Publish(bin)
		}
	}

	ctx.IndentedJSON(http.StatusOK, resp)
}
