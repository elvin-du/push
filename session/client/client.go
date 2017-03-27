/*
SESSION Server对外提供的客户端接口,本接口内部调用RPC服务
*/

package client

import (
	"hscore/log"
	"push/common/client/service"
	"push/meta"

	"golang.org/x/net/context"
)

func Online(req *meta.SessionOnlineRequest) (*meta.SessionOnlineResponse, error) {
	cli, err := service.SessionClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.SessionPut(cli)

	return cli.SessionClient.Online(context.TODO(), req)
}

func Offline(req *meta.SessionOfflineRequest) (*meta.SessionOfflineResponse, error) {
	cli, err := service.SessionClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.SessionPut(cli)

	return cli.SessionClient.Offline(context.TODO(), req)
}

func Info(req *meta.SessionInfoRequest) (*meta.SessionInfoResponse, error) {
	cli, err := service.SessionClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer service.SessionPut(cli)

	return cli.SessionClient.Info(context.TODO(), req)
}
