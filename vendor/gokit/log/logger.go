package log

import (
    "fmt"
    "github.com/vuleetu/logrus"
    "sync"
    "io"
    "runtime"
    "strings"
)

func New() *Logger {
    return &Logger{
        logrus.New(),
        logrus.DebugLevel,
        sync.Mutex{},
        false,
    }
}

type Logger struct {
    logger *logrus.Logger
    level logrus.Level
    mu sync.Mutex
    trace bool
}

type (
    Formatter logrus.Formatter
    Hook      logrus.Hook
)

var levelMapping = map[string]logrus.Level{
    "all":   logrus.DebugLevel,
    "debug": logrus.DebugLevel,
    "info":  logrus.InfoLevel,
    "warn":  logrus.WarnLevel,
    "error": logrus.ErrorLevel,
    "fatal": logrus.FatalLevel,
    "panic": logrus.PanicLevel,
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

func (l *Logger) SetStringLevel(level string) {
    ilevel, ok := levelMapping[level]
    if !ok {
        ilevel = logrus.DebugLevel
    }

    l.SetLevel(ilevel)
}

// SetOutput sets the standard logger output.
func (l *Logger) SetOutput(out io.Writer) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Out = out
}

// SetFormatter sets the standard logger formatter.
func (l *Logger) SetFormatter(formatter Formatter) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Formatter = formatter
}

// SetLevel sets the standard logger level.
func (l *Logger) SetLevel(level logrus.Level) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Level = level
}


func (l *Logger) AddHook(hook Hook) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Hooks.Add(hook)
}

func (l* Logger) SetTrace(trace bool) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.trace = trace
}

func (l *Logger) fileLine(skip int, args ...interface{}) []interface{} {
    if l.trace {
        _, file, line, _ := runtime.Caller(skip)
        fl := "["+fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/src/")+5:], line) + "]  "

        var arg []interface{}
        arg = append(arg, fl)
        arg = append(arg, args...)
        return arg
    }

    return args
}

func (l *Logger) Debug(args ...interface{}) {
    l.DebugN(3, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
    l.DebugfN(3, format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
    l.DebuglnN(3, args...)
}

func (l *Logger) Info(args ...interface{}) {
    l.InfoN(3, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
    l.InfofN(3, format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
    l.InfolnN(3, args...)
}

func (l *Logger) Warn(args ...interface{}) {
    l.WarnN(3, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
    l.WarnfN(3, format, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
    l.WarnlnN(3, args...)
}

func (l *Logger) Error(args ...interface{}) {
    l.ErrorN(3, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
    l.ErrorfN(3, format, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
    l.ErrorlnN(3, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
    l.FatalN(3, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
    l.FatalfN(3, format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
    l.FatallnN(3, args...)
}

func (l *Logger) Panic(args ...interface{}) {
    l.PanicN(3, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
    l.PanicfN(3, format, args...)
}

func (l *Logger) Panicln(args ...interface{}) {
    l.PaniclnN(3, args...)
}

//
func (l *Logger) DebugN(skip int, args ...interface{}) {
    l.logger.Debug(l.fileLine(skip, args...)...)
}

func (l *Logger) DebugfN(skip int, format string, args ...interface{}) {
    l.logger.Debugf(format, l.fileLine(skip, args...)...)
}

func (l *Logger) DebuglnN(skip int, args ...interface{}) {
    l.logger.Debugln(l.fileLine(skip, args...)...)
}

func (l *Logger) InfoN(skip int, args ...interface{}) {
    l.logger.Info(l.fileLine(skip, args...)...)
}

func (l *Logger) InfofN(skip int, format string, args ...interface{}) {
    l.logger.Infof(format, l.fileLine(skip, args...)...)
}

func (l *Logger) InfolnN(skip int, args ...interface{}) {
    l.logger.Infoln(l.fileLine(skip, args...)...)
}

func (l *Logger) WarnN(skip int, args ...interface{}) {
    l.logger.Warn(l.fileLine(skip, args...)...)
}

func (l *Logger) WarnfN(skip int, format string, args ...interface{}) {
    l.logger.Warnf(format, l.fileLine(skip, args...)...)
}

func (l *Logger) WarnlnN(skip int, args ...interface{}) {
    l.logger.Warnln(l.fileLine(skip, args...)...)
}

func (l *Logger) ErrorN(skip int, args ...interface{}) {
    l.logger.Error(l.fileLine(skip, args...)...)
}

func (l *Logger) ErrorfN(skip int, format string, args ...interface{}) {
    l.logger.Errorf(format, l.fileLine(skip, args...)...)
}

func (l *Logger) ErrorlnN(skip int, args ...interface{}) {
    l.logger.Errorln(l.fileLine(skip, args...)...)
}

func (l *Logger) FatalN(skip int, args ...interface{}) {
    l.logger.Fatal(l.fileLine(skip, args...)...)
}

func (l *Logger) FatalfN(skip int, format string, args ...interface{}) {
    l.logger.Fatalf(format, l.fileLine(skip, args...)...)
}

func (l *Logger) FatallnN(skip int, args ...interface{}) {
    l.logger.Fatalln(l.fileLine(skip, args...)...)
}

func (l *Logger) PanicN(skip int, args ...interface{}) {
    l.logger.Panic(l.fileLine(skip, args...)...)
}

func (l *Logger) PanicfN(skip int, format string, args ...interface{}) {
    l.logger.Panicf(format, l.fileLine(skip, args...)...)
}

func (l *Logger) PaniclnN(skip int, args ...interface{}) {
    l.logger.Panicln(l.fileLine(skip, args...)...)
}
