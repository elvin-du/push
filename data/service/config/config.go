package config

import (
	"hscore/config"
	"hscore/log"
	"hscore/util"
)

var (
	RpcServiceName    string
	RpcServiceVersion string
	RpcServicePort    string

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
	err := config.Get("service:rpc:name", &RpcServiceName)
	if nil != err {
		log.Fatal(err)
	}

	err = config.Get("service:rpc:version", &RpcServiceVersion)
	if nil != err {
		log.Fatal(err)
	}

	err = config.Get("service:rpc:port", &RpcServicePort)
	if nil != err {
		log.Fatal(err)
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
