package logging

import (
	"hot-coffee/config"
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger() error {
	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}

	logger = slog.New(slog.NewJSONHandler(logFile, nil))
	return nil
}

func Info(msg string, args ...interface{}) {
	if logger != nil {
		logger.Info(msg, args)
	}
}

func Error(msg string, err error, args ...interface{}) {
	if logger != nil {
		logger.Error(msg, append(args, "error", err.Error()))
	}
}

func Warn(msg string, args ...interface{}) {
	if logger != nil {
		logger.Warn(msg, args)
	}
}

func Fatal(msg string, err error, args ...interface{}) {
	if logger != nil {
		logger.Error(msg, append(args, "error", err.Error()))
		os.Exit(1)
	}
}
