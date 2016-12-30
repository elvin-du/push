package rpc

import (
	"log"
	"push/common/util"
	"push/meta"
)

const (
	MAX_POOL_SIZE  = 10
	SERVER_NAME    = "DATA"
	SERVER_VERSION = "1.0"
)

type DataRPCClientManager struct {
	Pool chan *DataRPCClient
}

type DataRPCClient struct {
	DataClient meta.DataClient
	RpcClient  *RPCClient
}

var (
	globalDataRPCClientManager = &DataRPCClientManager{Pool: make(chan *DataRPCClient, MAX_POOL_SIZE)}
)

func DataClient() (*DataRPCClient, error) {
	return globalDataRPCClientManager.GetClient()
}

func DataPut(cli *DataRPCClient) error {
	return globalDataRPCClientManager.Put(cli)
}

func (e *DataRPCClientManager) GetClient() (*DataRPCClient, error) {
	select {
	case cli := <-e.Pool:
		return cli, nil
	default:
		rpcCli, err := NewRPCClient(util.APP_NAME, SERVER_NAME, SERVER_VERSION)
		if nil != err {
			log.Println(err)
			return nil, err
		}

		dataCli := meta.NewDataClient(rpcCli.Client)

		return &DataRPCClient{
			DataClient: dataCli,
			RpcClient:  rpcCli,
		}, nil
	}
}

func (e *DataRPCClientManager) Put(cli *DataRPCClient) error {
	//discard client when pool is full
	//TODO
	select {
	case e.Pool <- cli:
		return nil
	default:
		err := cli.RpcClient.Close()
		if nil != err {
			log.Println(err)
		}

		return err
	}
}
