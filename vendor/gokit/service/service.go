package service

import (
	"gokit/config"
	"gokit/log"
	"gokit/service/pub"
	"gokit/util"
	"os"
)

type Option struct {
	PubNames []string
}

func Start(opts ...Option) {
	loadConfig()
	initLog()
	//启动pulisher服务
	for _, o := range opts {
		for _, p := range o.PubNames {
			pub.CreatePublisher(p)
		}
	}
	startLocalServices()
}

func loadConfig() {
	err := config.ReadConfig(util.GetFile("config.yml"))
	if err != nil {
		log.Fatal("Read configuration file failed", err)
	}
}

func initLog() {
	var enable_line bool
	err := config.Get("log:enable_line", &enable_line)
	if err != nil {
		log.Fatalln("Can not find log:enable_line, err: ", err)
	}

	log.EnableFileLine(enable_line)

	//set level
	var level string
	err = config.Get("log:level", &level)
	if err != nil {
		level = "all"
	}
	log.SetStringLevel(level)

	//set output
	var logoutput string
	err = config.Get("log:output", &logoutput)
	if err == nil && logoutput == "file" {
		var filename string
		err := config.Get("log:name", &filename)
		if err != nil {
			log.Fatal("Pleasde speficy file name when using file as log output")
		}

		filepath := util.GetFile(filename)
		f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal("Create log file failed")
		}

		log.SetOutput(f)
	}
}

func startLocalServices() {
	for _, srv := range _localServices {
		srv()
	}
}

var _localServices = make([]func(), 0, 10)

func Register(f func()) {
	_localServices = append(_localServices, f)
}
