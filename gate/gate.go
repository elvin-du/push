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
	resp := &meta.GatePushResponse{}

	svc := defaultServer.Services[req.Header.AppId+req.ClientId]
	err := svc.Push(uint16(req.PacketId), []byte(req.Content))
	if nil != err {
		return resp, err
	}

	return resp, nil
}

func (*Gate) PushAll(ctx context.Context, req *meta.GatePushAllRequest) (*meta.GatePushAllResponse, error) {
	log.Debugln("content", req.Content)
	return &meta.GatePushAllResponse{}, nil
}
