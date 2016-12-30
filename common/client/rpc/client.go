package rpc

import (
	//	"golang.org/x/net/context"
	"log"
	"push/common/client/etcd"

	"google.golang.org/grpc"
)

type RPCClient struct {
	Client        *grpc.ClientConn
	AppName       string
	ServerName    string
	ServerVersion string
}

func NewRPCClient(appName, srvName, srvVer string) (*RPCClient, error) {
	etcdCli, err := etcd.GetClient()
	if nil != err {
		log.Println(err)
		return nil, err
	}

	ip, port, err := etcdCli.Get(appName, srvName, srvVer)
	if nil != err {
		log.Println(err)
		return nil, err
	}

	target := ip + port
	cliConn, err := grpc.Dial(target, grpc.WithInsecure())
	if nil != err {
		log.Println(err)
		return nil, err
	}

	return &RPCClient{
		Client:        cliConn,
		AppName:       appName,
		ServerName:    srvName,
		ServerVersion: srvVer,
	}, nil
}

func (r *RPCClient) Close() error {
	return r.Client.Close()
}
