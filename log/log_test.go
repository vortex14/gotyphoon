package log

import (
	"errors"
	"testing"
)

func TestErrorLog(t *testing.T) {
	InitD()
	l := New(map[string]interface{}{"test": "test-logger"})
	l.Error("My error")
	l.Errorf("e: %v", errors.New("a new exception"))
}
