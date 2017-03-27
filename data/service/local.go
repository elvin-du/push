package service

import (
	"push/common/log"
	"push/data/service/config"
)

func Start() {
	config.Init()
	log.Init()
}
