package dal

import (
	"push/meta"
)

type Redis struct{}

//上线
func (r *Redis) Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	return &meta.DataOnlineResponse{}, nil
}

//下线
func (r *Redis) Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	return &meta.DataOfflineResponse{}, nil
}

func (r *Redis) SaveOfflineMsg(req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error) {
	return &meta.SaveOfflineMsgResponse{}, nil
}

func (r *Redis) GetOfflineMsgs(req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error) {
	return &meta.GetOfflineMsgsResponse{}, nil
}

func (r *Redis) DelOfflineMsgs(req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error) {
	return &meta.DelOfflineMsgsResponse{}, nil
}

func (r *Redis) GetClientInfo(req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error) {
	return &meta.GetClientInfoResponse{}, nil
}

func (r *Redis) UpdateClientInfo(req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error) {
	return &meta.UpdateClientInfoResponse{}, nil
}
