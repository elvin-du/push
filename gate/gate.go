/*
GATE对外提供RPC服务接口
*/

package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	"push/pb"

	"push/gate/message"

	"push/common/model"

	"golang.org/x/net/context"
)

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	AppID   string `json:"app_id"`
	Kind    int    `json:"kind"`
}

type Gate struct {
}

func (*Gate) Push(ctx context.Context, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	log.Debugln(*req)

	offlineMsg := &model.OfflineMsg{}
	offlineMsg.AppID = req.AppID
	offlineMsg.RegID = req.RegID
	offlineMsg.Content = req.Content
	offlineMsg.Extra = req.Extra
	offlineMsg.Kind = req.Kind
	offlineMsg.ID = req.ID
	message.DefaultMessageManager.Put(offlineMsg)
	var (
		err error = nil
		bin []byte
	)

	defer func() {
		if nil != err {
			message.DefaultMessageManager.Delete(offlineMsg.ID)
			err = model.OfflineMsgModel().Insert(offlineMsg)
			if nil != err {
				log.Errorln(err)
			}
		}
	}()

	user := defaultServer.Get(req.AppID, req.RegID)
	if nil == user {
		log.Errorln("not found session by:appID:regID:", req.AppID, req.RegID)
		return nil, errors.New("not found")
	}

	msg := Message{}
	msg.Content = req.Content
	msg.ID = req.ID
	msg.AppID = req.AppID
	msg.Kind = int(req.Kind)
	bin, err = json.Marshal(msg)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	err = user.Push(bin)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	resp := &pb.GatePushResponse{}
	return resp, nil
}

func (*Gate) PushAll(ctx context.Context, req *pb.GatePushAllRequest) (*pb.GatePushAllResponse, error) {
	log.Debugln("content", req.Content)
	return &pb.GatePushAllResponse{}, nil
}
