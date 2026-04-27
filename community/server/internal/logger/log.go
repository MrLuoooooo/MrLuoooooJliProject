package logger

import (
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(c *Config) {
	writeSyncer := getLogWriter(c)                                                 // 1. 获取日志写入目标（同时写入文件+控制台）
	encoder := getEncoder()                                                        // 2. 获取日志编码器（定义日志输出格式）
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)              // 3. 创建Zap日志核心（编码器+写入目标+日志级别)
	logger := zap.New(core, zap.AddStacktrace(zapcore.WarnLevel), zap.AddCaller()) // 4. 创建日志器实例（添加堆栈跟踪+调用者信息）
	zap.ReplaceGlobals(logger)                                                     // 5. 将该日志器设为全局默认日志器，后续可通过zap.L()/zap.S()直接调用
}

// getEncoder creates and returns a configured encoder for the logger.
func getEncoder() zapcore.Encoder {
	// 1. 基于Zap生产环境编码器配置为基础（自带默认的字段命名和格式）
	encoderConfig := zap.NewProductionEncoderConfig()
	// 2. 配置时间字段：使用ISO8601格式（如2024-05-20T15:30:59.123Z），字段名改为"time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	// 3. 配置级别字段：输出大写带颜色的级别（如INFO为蓝色、ERROR为红色），便于控制台区分
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 4. 配置耗时字段：以秒为单位输出（保留小数），便于直观查看操作耗时
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// 5. 配置调用者字段：输出短路径（如router/http.go:50），而非完整绝对路径，节省日志空间
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 6. 创建并返回控制台编码器（适配控制台输出，带颜色、分行等格式）
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter returns a WriteSyncer that writes logs to os.Stdout.
func getLogWriter(c *Config) zapcore.WriteSyncer {
	// 将io.Writer（多输出流）转换为zapcore.WriteSyncer（Zap支持的写入接口）
	return zapcore.AddSync(NewMultiWrite(c))
}

func NewMultiWrite(c *Config) io.Writer {
	// 1. 初始化lumberjack日志切割器，配置日志文件参数
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Logger.Path + c.Logger.FileName, // 日志文件完整路径（目录+文件名）
		MaxSize:    c.Logger.MaxSize,                  // 单个日志文件最大大小（单位：MB）
		MaxAge:     c.Logger.MaxAge,                   // 日志文件最大保存天数（超过自动删除）
		MaxBackups: c.Logger.MaxBackups,               // 最大保留备份文件数量（旧日志文件）
		Compress:   false,                             // 是否压缩备份文件（gzip）
	}

	// 2. 将lumberjack日志器转换为zapcore.WriteSyncer，再转为io.Writer
	syncFile := zapcore.AddSync(lumberJackLogger)
	// 3. 将标准输出（控制台）转换为zapcore.WriteSyncer，再转为io.Writer
	syncConsole := zapcore.AddSync(os.Stdout)
	// 4. 创建多输出流：写入该io.Writer的数据，会同时写入文件和控制台
	return io.MultiWriter(syncFile, syncConsole)
}
