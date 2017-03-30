package config

import (
	"hscore/config"
	"hscore/log"
	"hscore/util"
)

var (
	HTTP_ADDR    string
	HTTP_PROFILE bool
	HTTP_MODE    string
)

var (
	NSQ_ADDR string
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
	err := config.Get("http:addr", &HTTP_ADDR)
	if nil != err {
		log.Fatal(err)
	}

	err = config.Get("http:profile", &HTTP_PROFILE)
	if nil != err {
		log.Fatal(err)
	}

	err = config.Get("http:mode", &HTTP_MODE)
	if nil != err {
		log.Fatal(err)
	}

	err = config.Get("nsq:addr", &NSQ_ADDR)
	if nil != err {
		log.Fatal(err)
	}
}
