package logger

import (
	"flag"
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config  maxSize - в Mb, maxBackups, maxAge - в днях
type Config struct {
	Filename   string
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	StdOut     bool
}

// InitByParams создать глобальный логгер по параметрам конфига Config
func InitByParams(config Config) {
	loggerConfig := lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize, // megabytes
		MaxAge:     config.MaxAge,  // days
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	}
	Init(getCore(getAtomicLevel(config.Level), &loggerConfig, config.StdOut))
}

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

func getCore(level zap.AtomicLevel, loggerConfig *lumberjack.Logger, stdOut bool) zapcore.Core {
	var stdOutThread zapcore.WriteSyncer
	if stdOut {
		stdOutThread = zapcore.AddSync(os.Stdout)
	}

	file := zapcore.AddSync(loggerConfig)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	if stdOut {
		return zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdOutThread, level),
			zapcore.NewCore(fileEncoder, file, level),
		)
	}

	return zapcore.NewTee(
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel(logLevelValue string) zap.AtomicLevel {
	var level zapcore.Level
	var logLevel = flag.String("l", logLevelValue, "log level")
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
