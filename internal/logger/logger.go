package logger

import (
	configApp "github.com/alexMolokov/otus-rotate-banner/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zp *zap.SugaredLogger
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func New(c *configApp.LoggerConf) (*Logger, error) {
	cfg := zap.NewProductionConfig()

	cfg.Level = zap.NewAtomicLevelAt(getLevel(c.Level))
	cfg.Encoding = c.Encoding
	cfg.OutputPaths = []string{c.Output}
	cfg.ErrorOutputPaths = []string{c.Output}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02/Jan/2006:15:04:05 -0700")

	loggerZap, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		zp: loggerZap.Sugar(),
	}
	err = logger.zp.Sync()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.zp.Debugf(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.zp.Infof(msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.zp.Warnf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.zp.Errorf(msg, args...)
}
