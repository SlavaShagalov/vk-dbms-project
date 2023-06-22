package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func TestConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
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
	}
}

func NewTestLogger() (*zap.Logger, *os.File, error) {
	cfg := TestConfig()

	fileEncoder := zapcore.NewConsoleEncoder(cfg)
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	fileCore := zapcore.NewCore(fileEncoder, zapcore.Lock(zapcore.AddSync(file)), zapcore.DebugLevel)

	logger := zap.New(fileCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, file, nil
}
