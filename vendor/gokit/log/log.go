package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/vuleetu/logrus"
)

var (
	enableFileLine = false
)

func init() {
	SetFormatter(&TextFormatter{})
	SetStringLevel("all")
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

var levelMapping = map[string]logrus.Level{
	"all":   logrus.DebugLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
	"panic": logrus.PanicLevel,
}

func SetStringLevel(level string) {
	ilevel, ok := levelMapping[level]
	if !ok {
		ilevel = logrus.DebugLevel
	}

	SetLevel(ilevel)
}

func GetLevels(levelStr string) []logrus.Level {
	levels := make([]logrus.Level, 0, len(levelMapping))
	if levelStr == "" {
		Info("Sentry levels not speficied, will use all levels")
		for _, v := range levelMapping {
			levels = append(levels, v)
		}
	} else {
		levelArray := strings.Split(levelStr, "|")
		for _, l := range levelArray {
			if l != "" {
				if v, ok := levelMapping[l]; ok {
					levels = append(levels, v)
				}
			}
		}
	}

	return levels
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
	return logrus.GetLevel()
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

//Shortcut for level log
func Debug(args ...interface{}) {
	fileLine().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	fileLine().Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	fileLine().Debugln(args...)
}

func Info(args ...interface{}) {
	fileLine().Info(args...)
}

func Infof(format string, args ...interface{}) {
	fileLine().Infof(format, args...)
}

func Infoln(args ...interface{}) {
	fileLine().Infoln(args...)
}

func Warn(args ...interface{}) {
	fileLine().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	fileLine().Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	fileLine().Warnln(args...)
}

func Error(args ...interface{}) {
	fileLine().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	fileLine().Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	fileLine().Errorln(args...)
}

func Fatal(args ...interface{}) {
	fileLine().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	fileLine().Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	fileLine().Fatalln(args...)
}

func Panic(args ...interface{}) {
	fileLine().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	fileLine().Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	fileLine().Panicln(args...)
}

func fileLine() *logrus.Entry {
	feilds := make(logrus.Fields)
	if enableFileLine {
		_, file, line, _ := runtime.Caller(2)
		fl := fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/src/")+5:], line)
		return logrus.WithField("FileLine", fl).WithFields(feilds)
	}

	return logrus.NewEntry(logrus.StdLog()).WithFields(feilds)
}

func EnableFileLine(enable bool) {
	enableFileLine = enable
}

func StdLog() *logrus.Logger {
	return logrus.StdLog()
}
