package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New(level int8, logFilePath string) (*Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	cfg.DisableStacktrace = true
	cfg.OutputPaths = []string{
		logFilePath,
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("can't build logger: %w", err)
	}
	return &Logger{logger: logger}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) String(key string, val string) zap.Field {
	return zap.Field{Key: key, Type: zapcore.StringType, String: val} // nolint: exhaustivestruct
}

func (l *Logger) Int64(key string, val int64) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Int64Type, Integer: val} // nolint: exhaustivestruct
}

func (l *Logger) Duration(key string, val time.Duration) zap.Field {
	return zap.Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)} // nolint: exhaustivestruct
}
