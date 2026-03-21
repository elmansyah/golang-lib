package zap

import (
	"strings"
	
	"go.uber.org/zap/zapcore"
)

// StringToZapLevel converts a string representation of a log level
// into a zapcore.Level.
//
// Supported values (case-insensitive):
//   - "debug"
//   - "info"
//   - "warn", "warning"
//   - "error"
//   - "fatal"
//   - "panic"
//
// If the provided value is empty or unrecognized, the default level (def)
// will be returned instead.
//
func StringToZapLevel(value string, def zapcore.Level) zapcore.Level {
	switch strings.ToLower(value) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return def
	}
}
