package service

import (
	"errors"
	"gokit/log"
	"push/common/client/etcd"

	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client        *grpc.ClientConn
	ServerName    string
	ServerVersion string
}

func GetServiceClient(srvName, srvVer string) (*ServiceClient, error) {
	etcdCli, err := etcd.GetClient()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	ip, port, err := etcdCli.Get(srvName, srvVer)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return NewServiceClient(ip, port, srvName, srvVer)
}

func NewServiceClient(ip, port, srvName, srvVer string) (*ServiceClient, error) {
	if "" == ip || "" == port {
		err := errors.New("ip and port must not be empty")
		log.Errorln(err)
		return nil, err
	}
	target := ip + ":" + port
	cliConn, err := grpc.Dial(target, grpc.WithInsecure())
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	log.Infoln("NewServiceClient success,target:", target, "serviceName:", srvName, "serviceVersion:", srvVer)
	return &ServiceClient{
		Client:        cliConn,
		ServerName:    srvName,
		ServerVersion: srvVer,
	}, nil
}

func (r *ServiceClient) Close() error {
	return r.Client.Close()
}
