package config

import (
	"hscore/config"
	"hscore/log"
	"hscore/util"
)

var (
	NSQ_LOOKUPD_ADDRS = []string{}
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
	var tmp []interface{}
	err := config.Get("nsq:lookupd:addrs", &tmp)
	if nil != err {
		log.Fatal(err)
	}

	for _, val := range tmp {
		val2, ok := val.(string)
		if !ok {
			log.Fatalln(val, "cannot convert to string")
		}
		NSQ_LOOKUPD_ADDRS = append(NSQ_LOOKUPD_ADDRS, val2)
	}
}
