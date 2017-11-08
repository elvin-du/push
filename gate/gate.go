/*
GATE对外提供RPC服务接口
*/

package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	"push/pb"

	"golang.org/x/net/context"
)

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Gate struct {
}

func (*Gate) Push(ctx context.Context, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	log.Debugln(*req)
	resp := &pb.GatePushResponse{}

	user := defaultServer.Get(req.AppID, req.ClientId)
	if nil == user {
		log.Errorln("not found session by:appid:clientId:", req.AppID, req.ClientId)
		return nil, errors.New("not found")
	}

	msg := Message{}
	msg.Content = req.Content
	msg.ID = req.ID
	bin, err := json.Marshal(msg)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	err = user.Push(bin)
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
