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

	router.POST("/push", AuthHandler(_push.Push))

	go func() {
		internalRouter := gin.Default()
		internalRouter.GET("/apps", ListApp)
		internalRouter.POST("/app", CreateApp)
		internalRouter.PUT("/app", UpdateApp)
		internalRouter.DELETE("/app/:id", DeleteApp)
		internalRouter.Run(config.HTTP_INTERNAL_ADDR)
	}()

	go func() {
		appRouter := gin.Default()
		appRouter.POST("/register", AuthHandler(Register))
		appRouter.DELETE("/register/:id", AuthHandler(Unregister))
		appRouter.Run(config.HTTP_APP_ADDR)
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
