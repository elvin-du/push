package service

import (
	_ "push/gate/service/config"
	"push/gate/service/log"
)

func Start() {
	log.InitLog()
}
