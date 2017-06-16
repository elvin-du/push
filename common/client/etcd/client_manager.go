package etcd

import (
	"hscore/log"
	"time"
)

const (
	MAX_POOL_SIZE = 10
)

var (
	//etcd default config
	ECTD_ENDPOINTS  = []string{"127.0.0.1:2379"}
	DIAL_TIMEOUT    = time.Second * 10 //etcd dial timeout
	REQUEST_TIMEOUT = time.Second * 10 //etcd request timeout
)

type ETCDClientManager struct {
	Pool chan *ETCDClient
}

var (
	globalETCDClientManager = &ETCDClientManager{Pool: make(chan *ETCDClient, MAX_POOL_SIZE)}
)

func GetClient() (*ETCDClient, error) {
	return globalETCDClientManager.GetClient()
}

func Put(cli *ETCDClient) error {
	return globalETCDClientManager.Put(cli)
}

func (e *ETCDClientManager) GetClient() (*ETCDClient, error) {
	select {
	case cli := <-e.Pool:
		log.Infoln("got one etcd client from pool")
		return cli, nil
	default:
		log.Infoln("pool is empty,new one etcd client")
		cli, err := NewETCDClient(ECTD_ENDPOINTS, DIAL_TIMEOUT, REQUEST_TIMEOUT)
		if nil != err {
			log.Errorln(err)
			return nil, err
		}

		return cli, nil
	}
}

func (e *ETCDClientManager) Put(cli *ETCDClient) error {
	//discard client when pool is full
	//TODO
	select {
	case e.Pool <- cli:
		return nil
	default:
		err := cli.Close()
		if nil != err {
			log.Errorln(err)
		}

		return err
	}
}
