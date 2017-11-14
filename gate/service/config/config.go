package config

import (
	"gokit/config"
	"gokit/log"
	pushUtil "push/common/util"
)

var (
	SERVER_IP        string
	TCP_PORT         int
	RPC_SERVICE_PORT int
	AUTH_KEY         string
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

	err = config.Get("auth:key", &AUTH_KEY)
	if nil != err {
		log.Fatal(err)
	}
}
