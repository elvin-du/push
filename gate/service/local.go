package service

import (
	"gokit/service"
	"push/gate/service/config"
	"push/gate/service/db"
	"push/gate/service/session"
)

func Start() {
	service.Register(func() {
		config.Start()
		db.Start()
		session.Start()
	})

	service.Start()
}
