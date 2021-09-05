package log

import (
	"strings"
)

const (
	BASE = "BASE"
)

var DEFAULT = BaseOptions{
	ShowFile: true,
	FullTimestamp: true,
	ShowLine: true,
}

func GetLoggingConfig(name string, level string) BaseOptions {
	DEFAULT.Level = strings.ToUpper(level)

	switch name {
	case BASE:
		return DEFAULT
	}

	return DEFAULT
}