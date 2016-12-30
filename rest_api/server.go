package main

import (
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func StartHTTP() {
	router := gin.Default()
	router.GET("/push", _push.Push)
	router.Run(":60002")
}
