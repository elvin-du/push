package db

import (
	"fmt"
	"gokit/config"
	"gokit/log"
	dbredis "push/common/db"
)

var (
	redisPools = make(map[string]*dbredis.Pool)
)

func StartRedis() {
	startRedis("main")
}

func MainRedis() *dbredis.Pool {
	return redisPools["main"]
}

func startRedis(key string) {
	var (
		addr string
		pool int
	)

	err := config.Get(fmt.Sprintf("redis:%s:addr", key), &addr)
	if nil != err {
		log.Fatal(err)
	}
	err = config.Get(fmt.Sprintf("redis:%s:pool", key), &pool)
	if nil != err {
		log.Fatal(err)
	}

	opt := &dbredis.Option{
		MaxIdle:        dbredis.MAX_IDLE,
		MaxActive:      pool,
		IdleTimeout:    dbredis.IDLE_TIMEOUT,
		ConnectTimeout: dbredis.CONNECT_TIMEOUT,
		ReadTimeout:    dbredis.READ_TIMEOUT,
		WriteTimeout:   dbredis.WRITE_TIMEOUT,
	}
	redisPools[key] = dbredis.NewPoolWithOpt(addr, opt)
}
