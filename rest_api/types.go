package main

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*Context)

type Message struct {
	ID           string `json:"id"`
	RegID        string `json:"reg_id"`
	Platform     string `json:"platform"` //android,ios
	IsProduction bool   `json:"is_production"`
	Content      string `json:"content"`
	Kind         int    `json:"kind"`
	Extra        string `json:"extra"`
}

type Info struct {
	AppID string `json:"app_id"`
	*Message
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
