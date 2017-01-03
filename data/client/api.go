package client

import (
	"log"
	"push/common/client/service"
	"push/meta"

	"golang.org/x/net/context"
)

func Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Println(err)
		return nil, err
	}

	return cli.DataClient.Online(context.TODO(), req)
}

func Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Println(err)
		return nil, err
	}
	defer service.DataPut(cli)

	return cli.DataClient.Offline(context.TODO(), req)
}
