package main

import (
	"net/http"
	"net/http/pprof"
	"push/rest_api/service/config"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func StartHTTP() {
	gin.SetMode(config.HTTP_MODE)
	router = gin.Default()
	if config.HTTP_PROFILE {
		SetProfile()
	}

	router.GET("/push", _push.Push)

	router.Run(config.HTTP_ADDR)
}

func SetProfile() {
	router.GET("/debug/pprof/heap", Handler(pprof.Handler("heap")))
	router.GET("/debug/pprof/block", Handler(pprof.Handler("block")))
	router.GET("/debug/pprof/routine", Handler(pprof.Handler("goroutine")))
	router.GET("/debug/pprof/thread", Handler(pprof.Handler("threadcreate")))
	router.GET("/debug/pprof/mutex", Handler(pprof.Handler("mutex")))
}

func Handler(h http.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
