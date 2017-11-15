package main

import (
	"gokit/log"
	"push/common/model"
)

func Register(ctx *Context) {
	bin, err := ctx.GetRawData()
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	log.Debugf("Msg:%+v", string(bin))
	reg, err := ValidateRegisterReq(bin)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	registry, err := model.RegistryModel().Create(reg.AppID, reg.DevToken, reg.Kind)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	ctx.IndentedJSON(200, registry)
}

func Unregister(ctx *Context) {
	regID := ctx.Param("id")
	err := model.RegistryModel().Delete(regID)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}
}
