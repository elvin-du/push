package service

import (
	"push/common/db"
	"push/common/log"
	"push/gate/service/config"
)

func Start() {
	config.Init()
	log.Init()
	db.Init()
}
