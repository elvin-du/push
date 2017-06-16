package service

import (
	"hscore/log"
	"push/meta"
)

const (
	DATA_MAX_POOL_SIZE  = 10
	DATA_SERVER_NAME    = "DATA"
	DATA_SERVER_VERSION = "1.0.0"
)

type DataServiceClientManager struct {
	Pool chan *DataServiceClient
}

type DataServiceClient struct {
	DataClient    meta.DataClient
	ServiceClient *ServiceClient
}

var (
	globalDataServiceClientManager = &DataServiceClientManager{Pool: make(chan *DataServiceClient, DATA_MAX_POOL_SIZE)}
)

func DataClient() (*DataServiceClient, error) {
	return globalDataServiceClientManager.GetClient()
}

func DataPut(cli *DataServiceClient) error {
	return globalDataServiceClientManager.Put(cli)
}

func (e *DataServiceClientManager) GetClient() (*DataServiceClient, error) {
	select {
	case cli := <-e.Pool:
		return cli, nil
	default:
		srvCli, err := GetServiceClient(DATA_SERVER_NAME, DATA_SERVER_VERSION)
		if nil != err {
			log.Errorln(err)
			return nil, err
		}

		dataCli := meta.NewDataClient(srvCli.Client)

		return &DataServiceClient{
			DataClient:    dataCli,
			ServiceClient: srvCli,
		}, nil
	}
}

func (e *DataServiceClientManager) Put(cli *DataServiceClient) error {
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
