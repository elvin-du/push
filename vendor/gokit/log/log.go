package log

import (
    "github.com/vuleetu/logrus"
    "io"
)

var (
    std = New()
    logSrvs = map[string]*Logger{}
)

func Register(name string, log *Logger) {
    logSrvs[name] = log
}

func GetLogger(name string) *Logger {
    log, ok :=  logSrvs[name]
    if !ok {
        log = New()
        Register(name, log)
    }

    return log
}

func init() {
    SetFormatter(&TextFormatter{})
    SetStringLevel("all")
}

func SetStringLevel(level string) {
    std.SetStringLevel(level)
}

func SetOutput(out io.Writer) {
    std.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
    std.SetFormatter(formatter)
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
    std.SetLevel(level)
}
// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
    return std.logger.Level
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
    std.AddHook(hook)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
//func WithField(key string, value interface{}) *logrus.Entry {
    //return std.WithField(key, value)
//}

//// WithFields creates an entry from the standard logger and adds multiple
//// fields to it. This is simply a helper for `WithField`, invoking it
//// once for each field.
////
//// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
//// or Panic on the Entry it returns.
//func WithFields(fields logrus.Fields) *logrus.Entry {
    //return std.WithFields(fields)
//}

//Shortcut for level log
func Debug(args ...interface{}) {
    std.DebugN(3, args...)
}

func Debugf(format string, args ...interface{}) {
    std.DebugfN(3, format, args...)
}

func Debugln(args ...interface{}) {
    std.DebuglnN(3, args...)
}

func Info(args ...interface{}) {
    std.InfoN(3, args...)
}

func Infof(format string, args ...interface{}) {
    std.InfofN(3, format, args...)
}

func Infoln(args ...interface{}) {
    std.InfolnN(3, args...)
}

func Warn(args ...interface{}) {
    std.WarnN(3, args...)
}

func Warnf(format string, args ...interface{}) {
    std.WarnfN(3, format, args...)
}

func Warnln(args ...interface{}) {
    std.WarnlnN(3, args...)
}

func Error(args ...interface{}) {
    std.ErrorN(3, args...)
}

func Errorf(format string, args ...interface{}) {
    std.ErrorfN(3, format, args...)
}

func Errorln(args ...interface{}) {
    std.ErrorlnN(3, args...)
}

func Fatal(args ...interface{}) {
    std.FatalN(3, args...)
}

func Fatalf(format string, args ...interface{}) {
    std.FatalfN(3, format, args...)
}

func Fatalln(args ...interface{}) {
    std.FatallnN(3, args...)
}

func Panic(args ...interface{}) {
    std.PanicN(3, args...)
}

func Panicf(format string, args ...interface{}) {
    std.PanicfN(3, format, args...)
}

func Panicln(args ...interface{}) {
    std.PaniclnN(3, args...)
}

//out of date, use set trace instead
func EnableFileLine(enable bool) {
    SetTrace(enable)
}

func SetTrace(trace bool) {
    std.SetTrace(trace)
}
