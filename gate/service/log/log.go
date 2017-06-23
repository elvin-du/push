package log

import (
	"github.com/sirupsen/logrus"
)

var LogrusEntry *logrus.Entry

func Start() {
	logger := logrus.New()
	LogrusEntry = logrus.NewEntry(logger)
}
