package client

import (
	"gokit/log"
	"sync"
	"time"

	"github.com/processout/grpc-go-pool"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	pools = make(map[string]*grpcpool.Pool)
	mu    = &sync.RWMutex{}
)

const (
	POOL_MAX          = 50
	POOL_INIT         = 10
	POOL_IDLE_TIMEOUT = time.Second * 60 * 30
)

func Get(ip, port string) (*grpcpool.ClientConn, error) {
	mu.Lock()
	defer mu.Unlock()

	var err error = nil
	pool, ok := pools[target(ip, port)]
	if !ok {
		pool, err = grpcpool.New(Factory(ip, port), POOL_INIT, POOL_MAX, POOL_IDLE_TIMEOUT)
		if nil != err {
			log.Errorln(err)
			return nil, err
		}
		pools[target(ip, port)] = pool
	}
	log.Infof("pool capacity:%d,available:%d", pool.Capacity(), pool.Available())

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return pool.Get(ctx)
}

func target(ip, port string) string {
	return ip + ":" + port
}

func Factory(ip, port string) func() (*grpc.ClientConn, error) {
	return func() (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := grpc.DialContext(ctx, target(ip, port), grpc.WithInsecure())
		if err != nil {
			cancel()
			return nil, err
		}

		return conn, nil
	}
}
