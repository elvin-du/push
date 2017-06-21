package config

import (
	goconfig "gokit/config"
	"gokit/log"
)

var (
	NSQ_LOOKUPD_ADDRS = []string{}
)

func Start() {
	ParseConfig()
}

func ParseConfig() {
	var tmp []interface{}
	err := goconfig.Get("nsq:lookupd:addrs", &tmp)
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
