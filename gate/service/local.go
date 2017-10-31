package service

import (
	"gokit/service"
	"push/common/db"
	"push/gate/service/config"
)

func Start() {
	service.Register(func() {
		config.Start()
		db.StartRedis([]string{"main"})
		db.StartMysql([]string{"main", "shard1", "shard2", "shard3", "shard4", "shard5", "shard6", "shard7", "shard8", "shard9"})
	})

	service.Start()
}
