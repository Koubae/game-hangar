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

	vars "github.com/koubae/game-hangar"
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

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
}

type AppLogger struct {
	*zap.Logger
}

func (l *AppLogger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *AppLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *AppLogger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *AppLogger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *AppLogger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

func (l *AppLogger) Panic(msg string, fields ...zap.Field) {
	l.Logger.Panic(msg, fields...)
}

func (l *AppLogger) logCloser(loggerTmp *AppLogger) {
	err := l.Logger.Sync()
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

func (l *AppLogger) LogCloser(loggerTmp Logger, z *zap.Logger) {
	err := z.Sync()
	if err == nil {
		return
	}
	if _, ok := errors.AsType[*fs.PathError](err); !ok {
		loggerTmp.Error(
			"Error while shutting down logger",
			zap.String("type", fmt.Sprintf("%T", err)),
			zap.Error(err),
		)
	}
}

var logger *AppLogger

func GetLogger() *AppLogger {
	if logger == nil {
		panic("Logger not initialized")
	}
	return logger
}

func CreateLogger(logLevel LogLevel, filePath string) *AppLogger {
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

	_zapLogger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	logger = &AppLogger{
		Logger: _zapLogger,
	}
	return logger
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
}

// filepath: logs/app.log
// filePathAbs: <root-dir>/logs/app.log i.e. => /home/user/project/logs/app.log
// filePathDir: /home/user/project/logs
func createFileLoggerWriter(filePath string) *lumberjack.Logger {
	logger := zap.L()

	filePathAbs := filepath.Join(vars.RootDir, filePath)
	filePathDir := filepath.Dir(filePathAbs)

	if _, err := os.Stat(filePathDir); os.IsNotExist(err) {
		if err := os.MkdirAll(filePathDir, 0o755); err != nil {
			logger.Panic(
				"failed to create log directory",
				zap.String("dir", filePathDir),
				zap.Error(err),
			)
		}
	}

	return &lumberjack.Logger{
		Filename:   filePathAbs,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}
}
