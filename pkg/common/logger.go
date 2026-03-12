package common

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var LogLevels = [4]LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

func CreateLogger(logLevel LogLevel, filePath string) (*zap.Logger, func(loggerTmp *zap.Logger, logger *zap.Logger)) {
	stdout := zapcore.Lock(zapcore.AddSync(os.Stdout))

	var fileWriter zapcore.WriteSyncer
	if filePath == "" {
		fileWriter = stdout
	} else {
		fileWriter = zapcore.AddSync(createFileLoggerWriter(filePath))
	}

	zapLevel := zapcore.InfoLevel
	switch strings.ToLower(string(logLevel)) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		log.Fatal("Logger level invalid, must be one of: DEBUG, INFO, WARN, or ERROR")
	}
	level := zap.NewAtomicLevelAt(zapLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = utcTimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.EncodeTime = utcTimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel)), logCloser
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
}

func createFileLoggerWriter(filePath string) *lumberjack.Logger {
	logger := zap.L()
	dir := filepath.Dir(filePath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			logger.Panic(fmt.Sprintf("failed to create log directory %q, error: %v", dir, err))
		}
	}

	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}
}

func logCloser(loggerTmp *zap.Logger, logger *zap.Logger) {
	err := logger.Sync()
	if err == nil {
		return
	}

	/*  Ignore some errors related to closing the logger.
	@bug:
		- https://github.com/uber-go/zap/issues/772
		- https://github.com/uber-go/zap/issues/328
	 Also, this seems to work:
	!errors.Is(err, syscall.EINVAL)
	*/
	if _, ok := errors.AsType[*fs.PathError](err); !ok {
		loggerTmp.Error(
			"Error while shutting down logger",
			zap.String("type", fmt.Sprintf("%T", err)),
			zap.Error(err),
		)
	}
}
