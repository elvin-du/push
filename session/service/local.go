package service

import (
	"push/common/log"
	"push/gate/service/config"
)

func Start() {
	config.Init()
	log.Init()
}
