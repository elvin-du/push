package service

import (
	"push/common/log"
	"push/session/service/config"
)

func Start() {
	config.Init()
	log.Init()
}
