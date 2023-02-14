package interfaces

import "github.com/sirupsen/logrus"

const (
	DEBUG   = "DEBUG"
	WARNING = "WARNING"
	INFO    = "INFO"
	ERROR   = "ERROR"
)

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	WithFields(fields logrus.Fields) *logrus.Entry
}
