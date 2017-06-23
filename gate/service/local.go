package service

import (
	"gokit/service"
	"push/gate/service/config"
	"push/gate/service/db"
	"push/gate/service/log"
	"push/gate/service/session"
)

func Start() {
	service.Register(func() {
		config.Start()
		log.Start()
		db.Start()
		session.Start()

	})

	service.Start()
}
