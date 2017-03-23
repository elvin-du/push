package log

import (
	"hscore/config"
	"hscore/log"
	"hscore/util"
	"os"
)

func Init() {
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
