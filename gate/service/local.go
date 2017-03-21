package service

import (
	_ "push/gate/service/config"
	"push/common/log"
)

func Start() {
	log.InitLog()
}
