package client

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type pool struct {
	Clients map[string]*grpc.ClientConn
}

var DefaultPool = &pool{}

func init() {
	DefaultPool.Clients = make(map[string]*grpc.ClientConn)
}

func (p *pool) Get(ip, port string) (*grpc.ClientConn, error) {
	key := ip + ":" + port
	client := p.Clients[key]
	if nil == client {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := grpc.DialContext(ctx, key, grpc.WithInsecure())
		if err != nil {
			cancel()
			return nil, err
		}
		p.Clients[ip+port] = conn

		return conn, nil
	}

	return client, nil
}
