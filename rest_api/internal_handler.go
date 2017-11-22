package main

import (
	"encoding/json"
	"gokit/log"
	"push/common/model"

	"github.com/gin-gonic/gin"
)

func ListApp(ctx *gin.Context) {
	apps, err := model.GetApps()
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(500, err)
		return
	}

	ctx.IndentedJSON(200, apps)
}

func CreateApp(ctx *gin.Context) {
	bin, err := ctx.GetRawData()
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	app, err := ValidateAppReq(bin)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	err = model.SetApp(app)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	//	err = model.ReloadAppCache()
	//	if nil != err {
	//		log.Errorln(err)
	//		ctx.AbortWithError(400, err)
	//		return
	//	}
}

func UpdateApp(ctx *gin.Context) {
	bin, err := ctx.GetRawData()
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	app := model.App{}
	err = json.Unmarshal(bin, &app)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	err = model.AppModel().Update(&app)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	err = model.SetApp2Cache(&app)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(500, err)
		return
	}

	//	err = model.ReloadAppCache()
	//	if nil != err {
	//		log.Errorln(err)
	//		ctx.AbortWithError(400, err)
	//		return
	//	}
}

func DeleteApp(ctx *gin.Context) {
	appID := ctx.Param("id")
	err := model.AppModel().Delete(appID)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	err = model.DeleteAppFromCache(appID)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(400, err)
		return
	}

	//	err = model.ReloadAppCache()
	//	if nil != err {
	//		log.Errorln(err)
	//		ctx.AbortWithError(400, err)
	//		return
	//	}
}
