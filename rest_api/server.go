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

type Context struct {
	*gin.Context
	AppID string `json:"app_id"`
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

type HandlerFunc func(*Context)

func StartHTTP() {
	gin.SetMode(config.HTTP_MODE)
	router = gin.Default()
	if config.HTTP_PROFILE {
		SetProfile()
	}

	router.POST("/push", AuthHandler(_push.Push))

	go func() {
		internalRouter := gin.Default()
		internalRouter.POST("/ios/cert", AddIOSCert)
		internalRouter.Run(config.HTTP_INTERNAL_ADDR)
	}()

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

func AuthHandler(h HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pCtx := NewContext(ctx)
		err := Auth(pCtx)
		if nil != err {
			return
		}

		h(pCtx)
	}
}
