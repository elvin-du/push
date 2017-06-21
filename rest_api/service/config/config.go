package config

import (
	goconfig "gokit/config"
	"gokit/log"
)

var (
	HTTP_ADDR    string
	HTTP_PROFILE bool
	HTTP_MODE    string
)

var (
	NSQD_ADDR string
)

func Start() {
	ParseConfig()
}

func ParseConfig() {
	err := goconfig.Get("http:addr", &HTTP_ADDR)
	if nil != err {
		log.Fatal(err)
	}

	err = goconfig.Get("http:profile", &HTTP_PROFILE)
	if nil != err {
		log.Fatal(err)
	}

	err = goconfig.Get("http:mode", &HTTP_MODE)
	if nil != err {
		log.Fatal(err)
	}

	err = goconfig.Get("nsqd:addr", &NSQD_ADDR)
	if nil != err {
		log.Fatal(err)
	}
}
