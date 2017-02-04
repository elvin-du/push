/*
GATE Server对外提供的客户端接口,本接口内部调用RPC服务
*/

package client

import (
	"hscore/log"
	"push/common/client/service"
	"push/meta"

	"golang.org/x/net/context"
)

func Push(ip, port string, req *meta.GatePushRequest) (*meta.GatePushResponse, error) {
	cli, err := service.GateClient(ip, port)
	if nil != err {
		log.Error(err)
		return nil, err
	}
	defer service.GatePut(cli)

	return cli.GateClient.Push(context.TODO(), req)
}

//TODO maybe remove later
func PushAll(req *meta.GatePushAllRequest) (*meta.GatePushAllResponse, error) {
	//	cli, err := service.GateClient()
	//	if nil != err {
	//		log.Error(err)
	//		return nil, err
	//	}
	//	defer service.GatePut(cli)

	//	return cli.GateClient.PushAll(context.TODO(), req)
	return &meta.GatePushAllResponse{}, nil
}
