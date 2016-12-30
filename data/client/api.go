package client

import (
	"log"
	"push/common/client/rpc"
	"push/meta"

	"golang.org/x/net/context"
)

func Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	cli, err := rpc.DataClient()
	if nil != err {
		log.Println(err)
		return nil, err
	}

	return cli.DataClient.Online(context.TODO(), req)
}

func Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	cli, err := rpc.DataClient()
	if nil != err {
		log.Println(err)
		return nil, err
	}
	defer rpc.DataPut(cli)

	return cli.DataClient.Offline(context.TODO(), req)
}
