package log

import (
    "os"
    "hscore/util"
    "hscore/config"
    "hscore/log"
)

var (
    EntryLog *log.Logger
)

func InitLog() {
    logger := log.New()
    logger.SetFormatter(&log.TextFormatter{})

    var enable_line bool
    err := config.Get("entry_log:enable_line", &enable_line)
    if err != nil {
        log.Fatalln("Can not find entry_log:enable_line, err: ", err)
    }

    logger.SetTrace(enable_line)

    //set level
    var level string
    err = config.Get("entry_log:level", &level)
    if err != nil {
        level = "all"
    }
    logger.SetStringLevel(level)

    //set output
    var logoutput string
    err = config.Get("entry_log:output", &logoutput)
    if err == nil && logoutput == "file" {
        var filename string
        err := config.Get("entry_log:name", &filename)
        if err != nil {
            log.Fatal("Pleasde speficy file name when using file as log output")
        }

        filepath := util.GetFile(filename)
        f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
        if err != nil {
            log.Fatal("Create log file failed")
        }

        logger.SetOutput(f)
    }

    EntryLog = logger
}

