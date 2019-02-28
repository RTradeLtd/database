package database

import (
	"testing"

	"go.uber.org/zap"
)

func TestZapLogger(t *testing.T) {
	var l = NewZapLogger(LogLevelInfo, zap.NewExample().Sugar())
	l.Print("hello", "world")
	l = NewZapLogger(LogLevelDebug, zap.NewExample().Sugar())
	l.Print("goodbye", "world")
}
