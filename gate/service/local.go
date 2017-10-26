package service

import (
	"gokit/service"
	"push/gate/service/config"
	"push/gate/service/db"
	"push/gate/service/log"
)

func Start() {
	service.Register(func() {
		config.Start()
		log.Start()
		db.Start()
	})

	service.Start()
}
