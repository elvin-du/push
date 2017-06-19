package service

import (
	"push/common/log"
	"push/gate/service/config"
	"push/gate/service/db"
	"push/gate/service/session"
)

func Start() {
	config.Init()
	log.Init()
	db.Init()
	session.Start()
}
