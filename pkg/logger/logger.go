package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/project5e/web3-blog/pkg/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var Logger *zap.Logger

func Init(filename string, maxSize, maxBackup, maxAge int, compress bool, logType, level string) error {
	writerSyncer := getLoggerWriter(filename, maxSize, maxBackup, maxAge, compress, logType)
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return errors.WithStack(err)
	}
	core := zapcore.NewCore(getEncoder(), writerSyncer, logLevel)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

	zap.ReplaceGlobals(Logger)
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "Logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if app.IsLocal() || app.IsDev() || app.IsTest() || app.IsProduction() {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLoggerWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	if logType == "daily" {
		logName := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logName)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
		Compress:   compress,
	}

	if app.IsLocal() || app.IsDev() || app.IsTest() {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		return zapcore.AddSync(lumberJackLogger)
	}
}

func InfoIf(err error) {
	if err != nil {
		Logger.Info("Error Occurred: ", zap.Error(err))
	}
}

func WarnIf(err error) {
	if err != nil {
		Logger.Warn("Error Occurred: ", zap.Error(err))
	}
}

func ErrorIf(err error) {
	if err != nil {
		Logger.Error("Error Occurred: ", zap.Error(err))
	}
}

func Debug(message string, fields ...zap.Field) {
	Logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	Logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	Logger.Error(message, fields...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Sugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Sugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Logger.Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Sugar().Errorf(template, args)
}
