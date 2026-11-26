package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	File *os.File
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(),
		File:   l.File,
	}
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context")
	}
	return log
}

func NewLogger(config Config) (*Logger, error) {
	logLvl := zap.NewAtomicLevel()

	if err := logLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0744); err != nil {
		return nil, fmt.Errorf("mk log dir: %w", err)
	}

	timeStamp := time.Now().Format("2006-01-02T15-04-05")

	filePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("log.%s", timeStamp),
	)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("open or create log file: %w", err)
	}

	zapConf := zap.NewDevelopmentEncoderConfig()
	zapConf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05")

	zapEncoder := zapcore.NewConsoleEncoder(zapConf)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, file, logLvl),
		zapcore.NewCore(zapEncoder, os.Stdout, logLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		File:   file,
	}, nil
}
