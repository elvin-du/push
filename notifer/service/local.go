package service

import (
	"gokit/service"
	"push/notifer/service/config"
)

func Start() {
	service.Register(func() {
		config.Start()
	})

	service.Start()
}
