package service

import (
	"gokit/service"
	"push/common/db"
	"push/rest_api/service/config"
	"push/rest_api/service/nsq"
)

func Start() {
	service.Register(func() {
		config.Start()
		nsq.Start()
		db.StartMysql([]string{"main", "shard1", "shard2", "shard3", "shard4", "shard5", "shard6", "shard7", "shard8", "shard9"})
	})

	service.Start()
}
