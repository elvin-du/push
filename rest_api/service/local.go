package service

import (
	"gokit/service"
	"push/common/db"
	"push/common/model"
	"push/rest_api/service/config"
	"push/rest_api/service/nsq"
)

func Start() {
	service.Register(func() {
		config.Start()
		nsq.Start()
		//		db.StartMysql([]string{"main"})
		db.StartRedis([]string{"main"})
		db.StartMongo([]string{"main"})
		model.InitAppCache()
	})

	service.Start()
}
