package logger

import (
	"flag"
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitByParams создать глобальный логгер по параметрам, maxSize - в Mb, maxBackups, maxAge - в днях
func InitByParams(fileName, level string, maxSize, maxBackups, maxAge int, compress bool) {
	loggerConfig := lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxBackups, // days
		MaxBackups: maxAge,
		Compress:   compress,
	}
	Init(getCore(getAtomicLevel(level), &loggerConfig))
}

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

func getCore(level zap.AtomicLevel, loggerConfig *lumberjack.Logger) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(loggerConfig)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
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
