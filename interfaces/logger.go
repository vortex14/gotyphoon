package interfaces

import "github.com/sirupsen/logrus"

const (
	DEBUG 	 = "DEBUG"
	WARNING  = "WARNING"
	INFO   	 = "INFO"
	ERROR    = "ERROR"
)

type LoggerInterface interface {
	Debug    (args ...interface{})
	Info     (args ...interface{})
	Warning  (args ...interface{})
	Error    (args ...interface{})
	WithFields (fields logrus.Fields) *logrus.Entry
}