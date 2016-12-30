package etcd

import (
	"context"

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type ETCDClient struct {
	Client         *clientv3.Client
	Endpoints      []string
	DialTimeout    time.Duration
	RequestTimeout time.Duration
}

func NewETCDClient(endpoints []string, dialTimeout, requestTimeout time.Duration) (*ETCDClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if nil != err {
		log.Println(err)
		return nil, err
	}

	return &ETCDClient{
		Client:         cli,
		Endpoints:      endpoints,
		DialTimeout:    dialTimeout,
		RequestTimeout: requestTimeout,
	}, nil
}

/*
ETCD存储格式
key格式：/app/srvName/version/ip/port
value参数：
    weigth：权重．值：1~10,数字越大，权重越重 TODO
    load：现在负载大小．还没有想好怎么做　TODO
value参数格式：weigth=2&load=
*/
func (e *ETCDClient) Register(app, srvName, version, ip, port string, meta map[string]string) error {
	key := fmt.Sprintf("/%s/%s/%s/%s/%s", app, srvName, version, ip, port)
	value := ""
	for k, v := range meta {
		value += k + "=" + v + "&"
	}
	value = strings.TrimSuffix(value, "&")

	ctx, cancel := context.WithTimeout(context.Background(), e.RequestTimeout)
	//租赁１２０秒过期
	resp, err := e.Client.Grant(ctx, 120)
	cancel()
	if nil != err {
		log.Println(err)
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), e.RequestTimeout)
	_, err = e.Client.Put(ctx, key, value, clientv3.WithLease(resp.ID))
	cancel()
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}

func (e *ETCDClient) Heartbeat(app, srvName, version, ip, port string, meta map[string]string) error {
	return e.Register(app, srvName, version, ip, port, meta)
}

/*
获取经过均衡负载的服务的ip和port
*/
func (e *ETCDClient) Get(app, srvName, version string) (ip string, port string, err error) {
	key := fmt.Sprintf("/%s/%s/%s", app, srvName, version)
	ctx, cancel := context.WithTimeout(context.Background(), e.RequestTimeout)
	resp, err := e.Client.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if nil != err {
		log.Println(err)
		return "", "", err
	}

	//TODO 均衡负载
	if len(resp.Kvs) > 0 {
		key := resp.Kvs[0].Key
		//        val := resp.Kvs[0].Value
		log.Println("resp", string(key))
		keys := strings.Split(string(key), "/")
		log.Println(keys)
		ip = keys[4]
		port = keys[5]
	}

	return
}

func (e *ETCDClient) Close() error {
	return e.Client.Close()
}
