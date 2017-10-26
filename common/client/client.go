package client

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Get(ip, port string) (*grpc.ClientConn, error) {
	key := ip + ":" + port
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, key, grpc.WithInsecure())
	if err != nil {
		cancel()
		return nil, err
	}

	return conn, nil
}
