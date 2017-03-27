/*
DATA Server对外提供的客户端接口,本接口内部调用RPC服务
*/

package client

import (
	"hscore/log"
	"push/common/client/service"
	"push/meta"

	"golang.org/x/net/context"
)

func Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.Online(context.TODO(), req)
}

func Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.Offline(context.TODO(), req)
}

func SaveOfflineMsg(req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.SaveOfflineMsg(context.TODO(), req)
}

func GetOfflineMsgs(req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.GetOfflineMsgs(context.TODO(), req)
}

func DelOfflineMsgs(req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.DelOfflineMsgs(context.TODO(), req)
}

func GetClientInfo(req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.GetClientInfo(context.TODO(), req)
}

func UpdateClientInfo(req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.UpdateClientInfo(context.TODO(), req)
}
