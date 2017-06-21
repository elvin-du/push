package service

import (
	"gokit/service"
	"push/rest_api/service/config"
	"push/rest_api/service/nsq"
)

func Start() {
	service.Register(func() {
		config.Start()
		nsq.Start()
	})

	service.Start()
}
