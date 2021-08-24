package logger

import (
	"strings"

	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	BASE = "BASE"
)

var DEFAULT = interfaces.BaseLoggerOptions{
	ShowFile: true,
	FullTimestamp: true,
	ShowLine: true,
}

func GetLoggingConfig(name string, level string) interfaces.BaseLoggerOptions {
	DEFAULT.Level = strings.ToUpper(level)

	switch name {
	case BASE:
		return DEFAULT
	}

	return DEFAULT
}