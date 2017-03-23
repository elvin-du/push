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
	TcpPort           string
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

	err = config.Get("service:tcp:port", &TcpPort)
	if nil != err {
		log.Fatal(err)
	}
}
