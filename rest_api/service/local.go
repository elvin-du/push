package service

import (
	"push/common/log"
	"push/rest_api/service/config"
	"push/rest_api/service/nsq"
)

func Start() {
	config.Init()
	log.Init()
	nsq.Init()
}
