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
}

type BaseLoggerOptions struct {
	Name string
	Level string
	ShowLine bool
	ShowFile bool
	ShortFileName bool
	FullTimestamp bool
	level logrus.Level
}

func (o *BaseLoggerOptions) GetLevel(name string) logrus.Level {
	var level logrus.Level
	switch name {
	case DEBUG:
		level = logrus.DebugLevel
	case INFO:
		level = logrus.InfoLevel
	case ERROR:
		level = logrus.ErrorLevel
	case WARNING:
		level = logrus.WarnLevel
	}
	return level
}