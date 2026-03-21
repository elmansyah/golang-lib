package zap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup initializes a zap logger with file-based and console outputs.
//
// It builds multiple cores based on the provided log file configurations
// (e.g. success log, error log), each with its own level range.
//
// Behavior:
//   - File logs are always written in JSON format (structured logging).
//   - Console output is enabled for stdout with different verbosity:
//     * dev  → debug level
//     * prod → info level
//   - Each file receives logs only within its defined MinLevel–MaxLevel range.
//
// Notes:
//   - Caller should call logger.Sync() before application exit.
//   - Invalid level ranges (MinLevel > MaxLevel) will return an error.
func Setup(req *Params) (*zap.SugaredLogger, error) {
	logPath, err := validateRequest(req)
	if err != nil {
		return nil, err
	}
	
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	
	// normalize app mode
	mode := strings.ToLower(req.AppMode)
	if mode != "dev" && mode != "prod" {
		mode = "prod"
	}
	
	cores := make([]zapcore.Core, 0, len(req.LogFiles)+1)
	
	// build file cores
	for _, logRequest := range req.LogFiles {
		if logRequest.FileName == "" {
			return nil, errFileNameRequired
		}
		
		core, err := buildZap(logPath, req, logRequest, encoderCfg)
		if err != nil {
			return nil, err
		}
		
		cores = append(cores, core)
	}
	
	// console level based on mode
	var consoleLevel zapcore.Level
	if mode == "dev" {
		consoleLevel = zapcore.DebugLevel
	} else {
		consoleLevel = zapcore.InfoLevel
	}
	
	// console core (stdout only)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		consoleLevel,
	)
	
	cores = append(cores, consoleCore)
	
	combined := zapcore.NewTee(cores...)
	
	logger := zap.New(
		combined,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	
	return logger.Sugar(), nil
}

// buildZap creates a zapcore.Core for a single log file.
//
// It ensures the log file exists with proper permissions, then configures
// a rotating file writer using lumberjack.
//
// The core uses JSON encoding and filters logs based on the file's level range.
//
// Returns an error if:
//   - file cannot be created/opened
//   - MinLevel is greater than MaxLevel
func buildZap(logPath string, request *Params, file File, encoderCfg zapcore.EncoderConfig) (zapcore.Core, error) {
	fullPath := filepath.Join(logPath, file.FileName)
	
	// ensure file can be created with correct permission
	openFile, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, request.FilePermission)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errOpenLogFile, fullPath)
	}
	
	err = openFile.Close()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCloseLogFile, fullPath)
	}
	
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    request.MaxSize,
		MaxBackups: request.MaxBackups,
		MaxAge:     request.MaxAge,
		Compress:   request.Compress,
	})
	
	// file encoder always JSON (structured logging)
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	
	if file.MinLevel > file.MaxLevel {
		return nil, fmt.Errorf(
			"%w: file=%s min=%s max=%s",
			errInvalidLevelRange,
			file.FileName,
			file.MinLevel.String(),
			file.MaxLevel.String(),
		)
	}
	
	return zapcore.NewCore(
		encoder,
		writer,
		levelEnabler(file),
	), nil
}

// validateRequest validates the logger configuration and ensures
// the log directory exists (or creates it if missing).
//
// It returns the resolved log directory path or an error if validation fails.
func validateRequest(request *Params) (string, error) {
	if err := validateLogDir(request); err != nil {
		return "", err
	}
	
	if err := os.MkdirAll(request.LogDir, request.DirPermission); err != nil {
		return "", fmt.Errorf("%w: %w", errCreateDir, err)
	}
	
	return request.LogDir, nil
}

// validateLogDir performs basic validation on the request:
//
//   - request must not be nil
//   - at least one log file must be defined
//
// It also applies default permissions if not set.
func validateLogDir(request *Params) error {
	if request == nil {
		return errNilRequest
	}
	
	if len(request.LogFiles) == 0 {
		return errDirRequired
	}
	
	setDefaultPermission(request)
	
	return nil
}

// setDefaultPermission sets default directory and file permissions
// if they are not explicitly provided.
//
// Defaults:
//   - DirPermission  → 0755
//   - FilePermission → 0644
func setDefaultPermission(request *Params) {
	if request.DirPermission == 0 {
		request.DirPermission = 0755
	}
	
	if request.FilePermission == 0 {
		request.FilePermission = 0644
	}
}

// levelEnabler returns a LevelEnablerFunc that filters log entries
// based on the provided file's MinLevel and MaxLevel.
//
// Only logs within the inclusive range [MinLevel, MaxLevel] are written.
func levelEnabler(file File) zap.LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level >= file.MinLevel && level <= file.MaxLevel
	}
}
