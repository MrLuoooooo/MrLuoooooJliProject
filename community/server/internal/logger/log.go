package logger

import (
	"io"
	"os"

	"community-server/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	writeSyncer := getLogWriter(&cfg.Logger)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddStacktrace(zapcore.WarnLevel), zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return logger
}

func InitLogger(c *config.Config) {
	writeSyncer := getLogWriter(&c.Logger)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddStacktrace(zapcore.WarnLevel), zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(c *config.Logger) zapcore.WriteSyncer {
	return zapcore.AddSync(NewMultiWrite(c))
}

func NewMultiWrite(c *config.Logger) io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Path + c.FileName,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		Compress:   false,
	}

	syncFile := zapcore.AddSync(lumberJackLogger)
	syncConsole := zapcore.AddSync(os.Stdout)
	return io.MultiWriter(syncFile, syncConsole)
}
