package service

import (
	"gokit/service"
	"push/gate/service/config"
	"push/gate/service/db"
	"push/gate/service/log"
)

func Start() {
	service.Register(func() {
		log.Start()
		config.Start()
		db.StartMysql()
		db.StartRedis()
	})

	service.Start()
}
