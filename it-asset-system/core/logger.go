package core

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() {
	config := AppConfig.Log

	// Create lumberjack logger for info and error
	infoLumberjack := &lumberjack.Logger{
		Filename:   "logs/info.log",
		MaxSize:    config.MaxSize, // MB
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
		Compress:   true,
	}

	errorLumberjack := &lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// info level core
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	infoCore := zapcore.NewCore(encoder, zapcore.AddSync(infoLumberjack), infoLevel)

	// error level core
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorCore := zapcore.NewCore(encoder, zapcore.AddSync(errorLumberjack), errorLevel)

	// stdout core
	stdoutCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	core := zapcore.NewTee(infoCore, errorCore, stdoutCore)
	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}
