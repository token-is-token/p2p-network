package utils

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger struct {
	logger *zap.SugaredLogger
	level  LogLevel
	mu     sync.RWMutex
}

func NewLogger(name string, level LogLevel) (*Logger, error) {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(levelToZap(level)),
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: logger.Named(name).Sugar(),
		level:  level,
	}, nil
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *Logger) With(keysAndValues ...interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(keysAndValues...).Sugar(),
		level:  l.level,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func levelToZap(level LogLevel) zapcore.Level {
	switch level {
	case LogLevelDebug:
		return zapcore.DebugLevel
	case LogLevelInfo:
		return zapcore.InfoLevel
	case LogLevelWarn:
		return zapcore.WarnLevel
	case LogLevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func parseLevel(levelStr string) LogLevel {
	switch strings.ToLower(levelStr) {
	case "debug", "dbg":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error", "err":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}

func getEnvLevel() LogLevel {
	envLevel := os.Getenv("LOG_LEVEL")
	if envLevel == "" {
		return LogLevelInfo
	}
	return parseLevel(envLevel)
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func GetDefaultLogger() *Logger {
	once.Do(func() {
		logger, _ := NewLogger("p2p-network", getEnvLevel())
		defaultLogger = logger
	})
	return defaultLogger
}

func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}
