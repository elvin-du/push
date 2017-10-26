/*
GATE对外提供RPC服务接口
*/

package main

import (
	"errors"
	"gokit/log"
	"push/pb"

	"golang.org/x/net/context"
)

type Gate struct {
}

func (*Gate) Push(ctx context.Context, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	log.Debugln(*req)
	resp := &pb.GatePushResponse{}

	user := defaultServer.Get(req.AppID, req.ClientId)
	if nil == user {
		log.Errorln("not found service by:appid:clientId:", req.AppID, req.ClientId)
		return nil, errors.New("not found")
	}

	err := user.Push(uint16(req.PacketId), []byte(req.Content))
	if nil != err {
		log.Errorln(err)
		return resp, err
	}

	return resp, nil
}

func (*Gate) PushAll(ctx context.Context, req *pb.GatePushAllRequest) (*pb.GatePushAllResponse, error) {
	log.Debugln("content", req.Content)
	return &pb.GatePushAllResponse{}, nil
}
