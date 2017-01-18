package dal

import (
	"push/meta"
)

//数据接口层提供的接口定义
type DAL interface {
	Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error)
	Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error)
	SaveOfflineMsg(req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error)
	GetOfflineMsgs(req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error)
	DelOfflineMsgs(req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error)
	GetClientInfo(req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error)
	UpdateClientInfo(req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error)
}
