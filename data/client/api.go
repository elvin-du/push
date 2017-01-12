/*
DATA Server对外提供的客户端接口
*/

package client

import (
	"log"
	"push/common/client/service"
	"push/meta"

	"push/data/dal"

	"golang.org/x/net/context"
)

type Client struct {
	Dal      dal.DAL
	provider string //mysql,redis
}

func New(provider string) (*Client, error) {
	switch provider {
	case "mysql":
		return &Client{Dal: &dal.Mysql{}, provider: "mysql"}, nil
	case "redis":
		return &Client{Dal: &dal.Redis{}, provider: "redis"}, nil
	default:
	}

	//TODO
	return nil, nil
}

func (c *Client) SaveOfflineMsg() error {
	//TODO
	return c.Dal.SaveOfflineMsg()
}

func Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	cli, err := service.DataClient()
	if nil != err {
		log.Println(err)
		return nil, err
	}
	defer service.DataPut(cli)

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
