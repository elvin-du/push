package main

import (
	"push/data/dal"
	"push/meta"

	"golang.org/x/net/context"
)

type Data struct {
}

//上线
func (*Data) Online(ctx context.Context, req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).Online(req)
}

//下线
func (*Data) Offline(ctx context.Context, req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).Offline(req)
}

func (*Data) SaveOfflineMsg(ctx context.Context, req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).SaveOfflineMsg(req)
}

func (*Data) GetOfflineMsgs(ctx context.Context, req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).GetOfflineMsgs(req)
}

func (*Data) DelOfflineMsgs(ctx context.Context, req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).DelOfflineMsgs(req)
}

func (*Data) GetClientInfo(ctx context.Context, req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).GetClientInfo(req)
}

func (*Data) UpdateClientInfo(ctx context.Context, req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error) {
	return dal.DefaultClientManager.GetMysql(req.Header.AppName).UpdateClientInfo(req)
}
