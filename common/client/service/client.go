package service

import (
	//	"golang.org/x/net/context"
	"log"
	"push/common/client/etcd"

	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client        *grpc.ClientConn
	AppName       string
	ServerName    string
	ServerVersion string
}

func GetServieClient(appName, srvName, srvVer string) (*ServiceClient, error) {
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

	return NewServieClient(ip, port, appName, srvName, srvName)
}

func NewServieClient(ip, port, appName, srvName, srvVer string) (*ServiceClient, error) {
	target := ip + port
	cliConn, err := grpc.Dial(target, grpc.WithInsecure())
	if nil != err {
		log.Println(err)
		return nil, err
	}

	return &ServiceClient{
		Client:        cliConn,
		AppName:       appName,
		ServerName:    srvName,
		ServerVersion: srvVer,
	}, nil
}

func (r *ServiceClient) Close() error {
	return r.Client.Close()
}
