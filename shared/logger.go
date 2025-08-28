package config

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger returns a singleton zap.Logger instance
func GetLogger() *zap.Logger {
	once.Do(func() {
		cfg := zap.NewProductionEncoderConfig()
		cfg.TimeKey = "time"
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder

		fileEncoder := zapcore.NewJSONEncoder(cfg)
		// Open file in append mode, create if not exists, do not truncate
		logFile := "logs/app.log"
		fileWriter := zapcore.AddSync(&lumberjackLogger{filename: logFile})

		core := zapcore.NewCore(fileEncoder, fileWriter, zapcore.InfoLevel)
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
	return logger
}

// lumberjackLogger is a minimal zapcore.WriteSyncer for appending to a file without rotation
type lumberjackLogger struct {
	filename string
}

func (l *lumberjackLogger) Write(p []byte) (n int, err error) {
	f, err := openLogFile(l.filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(p)
}

func (l *lumberjackLogger) Sync() error {
	return nil
}

func openLogFile(path string) (*os.File, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 0644: owner read/write, group/others read
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
