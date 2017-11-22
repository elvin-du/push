/*
GATE对外提供RPC服务接口
*/

package main

import (
	"encoding/json"
	"errors"
	"gokit/log"
	//	"gokit/util"
	//	"push/common/model"
	//	"push/gate/message"
	"push/pb"

	"golang.org/x/net/context"
)

type PublishMessage struct {
	ID      string                 `json:"id"`
	Content string                 `json:"content"`
	AppID   string                 `json:"app_id"`
	Extras  map[string]interface{} `json:"extras"`
}

type Gate struct {
}

func (*Gate) Push(ctx context.Context, req *pb.GatePushRequest) (resp *pb.GatePushResponse, err error) {
	log.Debugln(*req)
	resp = &pb.GatePushResponse{}
	log.Infof("process msg_id:%s,app_id:%s,reg_id:%s", req.ID, req.AppID, req.RegID)

	var extras map[string]interface{}
	err = json.Unmarshal([]byte(req.Extras), &extras)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	//	offlineMsg := &model.Message{}
	//	offlineMsg.AppID = req.AppID
	//	offlineMsg.RegID = req.RegID
	//	offlineMsg.Content = req.Content
	//	offlineMsg.Extras = extras
	//	offlineMsg.ID = req.ID
	//	offlineMsg.CreatedAt = util.Timestamp()
	//	offlineMsg.Status = 1
	//	offlineMsg.TTL = req.TTL
	//	message.DefaultMessageManager.Put(offlineMsg)

	//	defer func() {
	//		if nil != err {
	//			message.DefaultMessageManager.Delete(offlineMsg.ID)
	//			err = model.MessageModel().Insert(offlineMsg)
	//			if nil != err {
	//				log.Errorln(err)
	//			}
	//		}
	//	}()

	user := defaultServer.Get(req.AppID, req.RegID)
	if nil == user {
		log.Errorf("not found session by appID:%s,regID:%s", req.AppID, req.RegID)
		return nil, errors.New("not found")
	}

	var bin []byte
	msg := PublishMessage{}
	msg.Content = req.Content
	msg.ID = req.ID
	msg.AppID = req.AppID
	msg.Extras = extras
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
	log.Infof("send msg success, msg_id:%s,app_id:%s,reg_id:%s", req.ID, req.AppID, req.RegID)
	return resp, nil
}

func (*Gate) PushAll(ctx context.Context, req *pb.GatePushAllRequest) (*pb.GatePushAllResponse, error) {
	log.Debugln("content", req.Content)
	return &pb.GatePushAllResponse{}, nil
}
