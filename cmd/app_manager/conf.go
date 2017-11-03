package main

import (
	"fmt"
	"gokit/config"
	"gokit/log"
)

var targets = []string{}

func InitConf() {
	t := []interface{}{}
	err := config.Get("targets:addr", &t)
	if nil != err {
		log.Fatal(err)
	}
	for _, v := range t {
		targets = append(targets, fmt.Sprintf("%s", v))
	}

	log.Debugf("%+v", targets)
}
