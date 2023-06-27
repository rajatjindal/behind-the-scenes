package slack

import "github.com/sirupsen/logrus"

type DebugLogging struct{}

func (d DebugLogging) Debug() bool {
	return true
}

func (d DebugLogging) Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func (d DebugLogging) Debugln(v ...interface{}) {
	logrus.Debug(v...)
}
