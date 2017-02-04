/*
GATE对外提供RPC服务接口
*/

package main

import (
	"hscore/log"
	"push/meta"

	"golang.org/x/net/context"
)

type Gate struct {
}

func (*Gate) Push(ctx context.Context, req *meta.GatePushRequest) (*meta.GatePushResponse, error) {
	log.Debugln(*req)
	svcs := defaultServer.Services[req.UserId]
	for _, v := range svcs {
		v.Push([]byte(req.Content))
	}
	return &meta.GatePushResponse{}, nil
}

func (*Gate) PushAll(ctx context.Context, req *meta.GatePushAllRequest) (*meta.GatePushAllResponse, error) {
	log.Debugln("content", req.Content)
	return &meta.GatePushAllResponse{}, nil
}
