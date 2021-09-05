package log

import "github.com/sirupsen/logrus"

const (
	DEBUG 	 = "DEBUG"
	WARNING  = "WARNING"
	INFO   	 = "INFO"
	ERROR    = "ERROR"
)

type BaseOptions struct {
	Name string
	Level string
	ShowLine bool
	ShowFile bool
	ShortFileName bool
	FullTimestamp bool
}

func (o *BaseOptions) GetLevel(name string) logrus.Level {
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