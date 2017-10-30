/*
GATE Server对外提供的客户端接口,本接口内部调用RPC服务
*/

package client

import (
	"gokit/log"
	"push/common/client"
	"push/pb"

	"golang.org/x/net/context"
)

func Push(ip, port string, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	return doPush(ip, port, req)
}

func doPush(ip, port string, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	cli, err := client.Get(ip, port)
	if nil != err {
		log.Error(err)
		return nil, err
	}
	defer cli.Close()

	gateCli := pb.NewGateClient(cli.ClientConn)

	return gateCli.Push(context.TODO(), req)
}

//TODO maybe remove later
func PushAll(req *pb.GatePushAllRequest) (*pb.GatePushAllResponse, error) {
	return &pb.GatePushAllResponse{}, nil
}
