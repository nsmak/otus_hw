package logger

import (
	"fmt"

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

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}
