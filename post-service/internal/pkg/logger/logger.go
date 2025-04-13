package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logFunc func(template string, args ...any)

type ILogger interface {
	Debug(args ...any)
	Debugf(template string, args ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
}

type Logger struct {
	logger *zap.SugaredLogger
	prefix string
	core   zapcore.Core
}

func NewLogger(name, prefix string) *Logger {
	logger, err := newSugaredLogger(name)
	if err != nil {
		return nil
	}
	log := wrap(logger, prefix)
	return log
}

func wrap(logger *zap.SugaredLogger, prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		logger: logger,
		core:   logger.Desugar().Core(),
	}
}

func newSugaredLogger(name string) (*zap.SugaredLogger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config := zap.NewProductionConfig()
	config.EncoderConfig = encoderConfig

	log, err := config.Build()
	if err != nil {
		return nil, err
	}
	log = log.Named(name)
	return log.Sugar(), nil
}

func getMessage(template string, args []any) string {
	if template != "" {
		return fmt.Sprintf(template, args...)
	}
	return fmt.Sprint(args...)
}

func (c *Logger) Debug(args ...any) {
	if !c.core.Enabled(zap.DebugLevel) {
		return
	}
	c.logger.Debugf(c.prefix+getMessage("", args), args...)
}

func (c *Logger) Debugf(template string, args ...any) {
	if !c.core.Enabled(zap.DebugLevel) {
		return
	}
	c.logger.Debugf(c.prefix+getMessage(template, args), args...)
}

func (c *Logger) Info(args ...any) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Infof(c.prefix+getMessage("", args), args...)
}

func (c *Logger) Infof(template string, args ...any) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Infof(c.prefix+getMessage(template, args), args...)
}

func (c *Logger) Error(args ...any) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Errorf(c.prefix+getMessage("", args), args...)
}

func (c *Logger) Errorf(template string, args ...any) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Errorf(c.prefix+getMessage(template, args), args...)
}
