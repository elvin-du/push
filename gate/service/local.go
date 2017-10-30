package service

import (
	"gokit/service"
	"push/gate/service/config"
	"push/gate/service/db"
)

func Start() {
	service.Register(func() {
		config.Start()
		db.StartMysql()
		db.StartRedis()
	})

	service.Start()
}
