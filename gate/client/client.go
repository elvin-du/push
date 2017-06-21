/*
GATE Server对外提供的客户端接口,本接口内部调用RPC服务
*/

package client

import (
	"gokit/log"
	"push/common/client/service"
	"push/pb"

	"golang.org/x/net/context"
)

func Push(ip, port string, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	//	ses, err := s.Get(req.ClientId)
	//	if nil != err {
	//		//TODO save push msg to offline msg
	//		log.Errorln(err)
	//		return nil, err
	//	}

	//	if "IOS" == ses.Platform {
	//		//TODO
	//		return nil, nil
	//	}
	//	//TODO req.PacketId需要填写

	return doPush(ip, port, req)
}

func doPush(ip, port string, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	cli, err := service.GateClient(ip, port)
	if nil != err {
		log.Error(err)
		return nil, err
	}
	defer service.GatePut(cli)

	return cli.GateClient.Push(context.TODO(), req)
}

//TODO maybe remove later
func PushAll(req *pb.GatePushAllRequest) (*pb.GatePushAllResponse, error) {
	return &pb.GatePushAllResponse{}, nil
}
