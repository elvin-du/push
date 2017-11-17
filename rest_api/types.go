package main

import (
	"push/common/types"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*Context)

type Message struct {
	*types.Message
}

func NewMessage() *Message {
	return &Message{
		&types.Message{},
	}
}

type Context struct {
	*gin.Context
	AppID string `json:"app_id"`
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{
		Context: ctx,
	}
}
