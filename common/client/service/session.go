package service

import (
	"hscore/log"
	"push/common/util"
	"push/meta"
)

const (
	SESSION_MAX_POOL_SIZE = 10
	SESSION_SERVER_NAME    = "SESSION"
	SESSION_SERVER_VERSION = "1.0.0"
)

type SessionServiceClientManager struct {
	Pool chan *SessionServiceClient
}

type SessionServiceClient struct {
	SessionClient meta.SessionClient
	ServiceClient *ServiceClient
}

var (
	globalSessionServiceClientManager = &SessionServiceClientManager{Pool: make(chan *SessionServiceClient, SESSION_MAX_POOL_SIZE)}
)

func SessionClient() (*SessionServiceClient, error) {
	return globalSessionServiceClientManager.GetClient()
}

func SessionPut(cli *SessionServiceClient) error {
	return globalSessionServiceClientManager.Put(cli)
}

func (e *SessionServiceClientManager) GetClient() (*SessionServiceClient, error) {
	select {
	case cli := <-e.Pool:
		log.Infoln("got one sessioin client from pool")
		return cli, nil
	default:
		log.Infoln("pool is empty,new one session client")
		srvCli, err := GetServiceClient(util.APP_NAME, SESSION_SERVER_NAME, SESSION_SERVER_VERSION)
		if nil != err {
			log.Errorln(err)
			return nil, err
		}

		sessionCli := meta.NewSessionClient(srvCli.Client)

		return &SessionServiceClient{
			SessionClient: sessionCli,
			ServiceClient: srvCli,
		}, nil
	}
}

func (e *SessionServiceClientManager) Put(cli *SessionServiceClient) error {
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
