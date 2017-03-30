package main

import (
	"hscore/log"
	dataCli "push/data/client"
	"push/meta"

	"golang.org/x/net/context"
)

type Session struct{}

func (*Session) Online(ctx context.Context, req *meta.SessionOnlineRequest) (*meta.SessionOnlineResponse, error) {
	log.Debugln("call Online,req:", *req)

	onlineReq := &meta.DataOnlineRequest{}
	onlineReq.ClientId = req.ClientId
	onlineReq.CreatedAt = req.CreatedAt
	onlineReq.GateServerIP = req.GateServerIP
	onlineReq.GateServerPort = req.GateServerPort
	onlineReq.Header = req.Header
	onlineReq.Platform = req.Platform
	resp, err := dataCli.Online(onlineReq)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.SessionOnlineResponse{
		Header: resp.Header,
	}, nil
}

func (*Session) Offline(ctx context.Context, req *meta.SessionOfflineRequest) (*meta.SessionOfflineResponse, error) {
	log.Debugln("call Offline,req: ", *req)

	offlineReq := &meta.DataOfflineRequest{}
	offlineReq.ClientId = req.ClientId
	offlineReq.Header = req.Header

	resp, err := dataCli.Offline(offlineReq)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.SessionOfflineResponse{
		Header: resp.Header,
	}, nil
}

func (*Session) Update(ctx context.Context, req *meta.SessionUpdateRequest) (*meta.SessionUpdateResponse, error) {
	log.Debugln("call Update,req: ", *req)

	updateReq := &meta.UpdateClientInfoRequest{}
	updateReq.ClientId = req.ClientId
	updateReq.GateServerIP = req.GateServerIP
	updateReq.GateServerPort = req.GateServerPort
	updateReq.Header = req.Header
	updateReq.Platform = req.Platform
	updateReq.UpdatedAt = req.UpdatedAt
	resp, err := dataCli.UpdateClientInfo(updateReq)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.SessionUpdateResponse{
		Header: resp.Header,
	}, nil
}

func (*Session) Info(ctx context.Context, req *meta.SessionInfoRequest) (*meta.SessionInfoResponse, error) {
	log.Debugln("call Info,req:", *req)

	infoReq := &meta.GetClientInfoRequest{}
	infoReq.ClientId = req.ClientId
	infoReq.Header = req.Header

	resp, err := dataCli.GetClientInfo(infoReq)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.SessionInfoResponse{
		Header:         resp.Header,
		AppId:          resp.AppId,
		ClientId:       resp.ClientId,
		CreatedAt:      resp.CreatedAt,
		GateServerIP:   resp.GateServerIP,
		GateServerPort: resp.GateServerPort,
		Platform:       resp.Platform,
		Status:         resp.Status,
		UpdatedAt:      resp.UpdatedAt,
	}, nil
}
