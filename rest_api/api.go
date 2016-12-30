package main

import (
	"log"
	"net/http"
	"push/data/client"
	"push/meta"

	"github.com/gin-gonic/gin"
)

type push struct{}

var _push *push

func (*push) Push(ctx *gin.Context) {
	resp, err := client.Online(&meta.DataOnlineRequest{UserId: "123", IP: "12.12.151.2"})
	if nil != err {
		log.Println(err)
	}
	log.Println(resp.String())

	ctx.String(http.StatusOK, "hi push")
}
