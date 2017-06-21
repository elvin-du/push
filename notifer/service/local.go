package service

import (
	"gokit/service"
	"push/notifer/service/config"
	"push/notifer/service/db"
)

func Start() {
	service.Register(func() {
		config.Start()
		db.StartRedis()
	})

	service.Start()
}
