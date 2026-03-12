package common

import (
	"log"
	"os"
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

func CreateLogger(logLevel LogLevel) *zap.Logger {
	stdout := zapcore.Lock(zapcore.AddSync(os.Stdout))

	file := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "/logs/app.log", // TODO: this should be a config too.-
			MaxSize:    10,              // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
		},
	)

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
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
}
