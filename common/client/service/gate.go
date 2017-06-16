package service

import (
	"hscore/log"
	"push/meta"
)

const (
	GATE_MAX_POOL_SIZE  = 10
	GATE_SERVER_NAME    = "GATE"
	GATE_SERVER_VERSION = "1.0.0"
)

type GateServiceClientManager struct {
	Pool chan *GateServiceClient
}

type GateServiceClient struct {
	GateClient    meta.GateClient
	ServiceClient *ServiceClient
}

var (
	globalGateServiceClientManager = &GateServiceClientManager{Pool: make(chan *GateServiceClient, GATE_MAX_POOL_SIZE)}
)

func GateClient(ip, port string) (*GateServiceClient, error) {
	return globalGateServiceClientManager.Client(ip, port)
}

func GatePut(cli *GateServiceClient) error {
	return globalGateServiceClientManager.Put(cli)
}

func (e *GateServiceClientManager) Client(ip, port string) (*GateServiceClient, error) {
	select {
	case cli := <-e.Pool:
		return cli, nil
	default:
		srvCli, err := NewServiceClient(ip, port, GATE_SERVER_NAME, GATE_SERVER_VERSION)
		if nil != err {
			log.Errorln(err)
			return nil, err
		}

		cli := meta.NewGateClient(srvCli.Client)

		return &GateServiceClient{
			GateClient:    cli,
			ServiceClient: srvCli,
		}, nil
	}
}

func (e *GateServiceClientManager) Put(cli *GateServiceClient) error {
	//discard client when pool is full
	//TODO
	select {
	case e.Pool <- cli:
		return nil
	default:
		err := cli.ServiceClient.Close()
		if nil != err {
			log.Errorln(err)
		}

		return err
	}
}
