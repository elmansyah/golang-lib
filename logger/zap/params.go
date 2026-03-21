package zap

import (
	"errors"
	"os"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	errNilRequest        = errors.New("logger config: request is nil")
	errDirRequired       = errors.New("logger config: base path is required")
	errFileNameRequired  = errors.New("logger config: file name is required")
	errCreateDir         = errors.New("logger config: failed to create directory")
	errOpenLogFile       = errors.New("logger config: failed to open log file")
	errCloseLogFile      = errors.New("logger config: failed to close log file")
	errInvalidLevelRange = errors.New("logger config: min level must be less than or equal to max level")
)

// Zap is a lightweight wrapper around zap.SugaredLogger.
//
// It provides a simplified logging interface while still allowing
// access to the underlying SugaredLogger for advanced use cases.
type Zap struct {
	Sugar *zap.SugaredLogger
}

// File represents a single log file configuration.
//
// Each file defines a log output target along with a level range.
// Only log entries within the inclusive range [MinLevel, MaxLevel]
// will be written to this file.
type File struct {
	FileName string
	MinLevel zapcore.Level
	MaxLevel zapcore.Level
}

// Params defines the configuration for initializing the logger.
//
// It includes settings for log file outputs, rotation, permissions,
// and environment mode.
//
// Fields:
//   - AppMode: application mode ("dev" or "prod"), affects console verbosity
//   - LogDir: directory where log files will be stored
//   - LogFiles: list of file configurations with level ranges
//   - DirPermission: permission for log directory (default: 0755)
//   - FilePermission: permission for log files (default: 0644)
//   - MaxSize: maximum size (MB) before log rotation
//   - MaxBackups: number of old log files to retain
//   - MaxAge: number of days to retain old log files
//   - Compress: whether to compress rotated log files
//
// This struct is typically constructed from environment variables.
type Params struct {
	AppMode        string
	LogDir         string
	LogFiles       []File
	DirPermission  os.FileMode
	FilePermission os.FileMode
	MaxSize        int
	MaxBackups     int
	MaxAge         int
	Compress       bool
}
