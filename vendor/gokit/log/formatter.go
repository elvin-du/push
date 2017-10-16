package log

import (
    "bytes"
    "fmt"
    "regexp"
    "sort"
    "strings"
    "time"
    "github.com/vuleetu/logrus"
)

const (
    nocolor = 0
    red     = 31
    green   = 32
    yellow  = 33
    blue    = 34
)

var (
    //baseTimestamp time.Time
    isTerminal    bool
    noQuoteNeeded *regexp.Regexp
)

func init() {
    //baseTimestamp = time.Now()
    isTerminal = logrus.IsTerminal()
}

type TextFormatter struct {
    // Set to true to bypass checking for a TTY before outputting colors.
    ForceColors   bool
    DisableColors bool
    // Set to true to disable timestamp logging (useful when the output
    // is redirected to a logging system already adding a timestamp)
    DisableTimestamp bool
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {

    var keys []string
    for k := range entry.Data {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    b := &bytes.Buffer{}

    //prefixFieldClashes(entry.Data)

    isColored := (f.ForceColors || isTerminal) && !f.DisableColors

    if isColored {
        printColored(b, entry, keys)
    } else {
        if !f.DisableTimestamp {
            f.appendKeyValue(b, "time", entry.Time.Format(time.RFC3339))
        }
        f.appendKeyValue(b, "level", entry.Level.String())
        f.appendKeyValue(b, "msg", entry.Message)
        for _, key := range keys {
            f.appendKeyValue(b, key, entry.Data[key])
        }
    }

    b.WriteByte('\n')
    return b.Bytes(), nil
}

func printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string) {
    var levelColor int
    switch entry.Level {
    case logrus.WarnLevel:
        levelColor = yellow
    case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
        levelColor = red
    default:
        levelColor = blue
    }

    levelText := strings.ToUpper(entry.Level.String())

    fmt.Fprintf(b, "%s \x1b[%dm[%s]\x1b[0m %-44s ", entry.Time.Format(time.RFC3339), levelColor, levelText, entry.Message)
    for _, k := range keys {
        v := entry.Data[k]
        fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
    }
}

func needsQuoting(text string) bool {
    for _, ch := range text {
        if !((ch >= 'a' && ch <= 'z') ||
            (ch >= 'A' && ch <= 'Z') ||
            (ch >= '0' && ch < '9') ||
            ch == '-' || ch == '.') {
            return false
        }
    }
    return true
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key, value interface{}) {
    switch value.(type) {
    case string:
        if needsQuoting(value.(string)) {
            fmt.Fprintf(b, "%v=%s ", key, value)
        } else {
            fmt.Fprintf(b, "%v=%q ", key, value)
        }
    case error:
        if needsQuoting(value.(error).Error()) {
            fmt.Fprintf(b, "%v=%s ", key, value)
        } else {
            fmt.Fprintf(b, "%v=%q ", key, value)
        }
    default:
        fmt.Fprintf(b, "%v=%v ", key, value)
    }
}

