package config

import (
	"gokit/config"
	"gokit/log"
	pushUtil "push/common/util"
)

var (
	RPC_SERVICE_PORT string
	SERVER_IP        string
	TCP_PORT         string
)

func Start() {
	ParseConfig()
}

func ParseConfig() {
	err := config.Get("service:rpc:port", &RPC_SERVICE_PORT)
	if nil != err {
		log.Fatal(err)
	}

	SERVER_IP = pushUtil.InternalIP

	err = config.Get("service:tcp:port", &TCP_PORT)
	if nil != err {
		log.Fatal(err)
	}
}
