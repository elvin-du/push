package config

import (
	"hscore/config"
	"hscore/log"
	"hscore/util"
	pushUtil "push/common/util"
)

var (
	RPC_SERVICE_PORT string

	SERVER_IP string

	MYSQL_DSN  string
	MYSQL_POOL int
)

func Init() {
	loadConfig()
	ParseConfig()
}

func loadConfig() {
	err := config.ReadConfig(util.GetFile("config.yml"))
	if err != nil {
		log.Fatal("Read configuration file failed", err)
	}
}

func ParseConfig() {
	err := config.Get("service:rpc:port", &RPC_SERVICE_PORT)
	if nil != err {
		log.Fatal(err)
	}

	externalIp := false
	err = config.Get("service:externalip", &externalIp)
	if nil != err {
		log.Fatal(err)
	}
	if externalIp {
		SERVER_IP = pushUtil.ExternalIP
	} else {
		SERVER_IP = pushUtil.InternalIP
	}

	err = config.Get("db:mysql:dsn", &MYSQL_DSN)
	if nil != err {
		log.Fatalln(err)
	}

	err = config.Get("db:mysql:pool", &MYSQL_POOL)
	if nil != err {
		log.Fatalln(err)
	}
}
