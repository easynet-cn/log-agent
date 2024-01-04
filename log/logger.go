package log

import (
	"os"
	"time"

	"github.com/easynet-cn/log-agent/util"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

var (
	Logger  *zap.Logger
	String  = zap.String
	Any     = zap.Any
	Int     = zap.Int
	Float32 = zap.Float32
)

func InitLogger(viper *viper.Viper) {
	hook := lumberjack.Logger{
		Filename:   viper.GetString("logging.file"),
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}
	write := zapcore.AddSync(&hook)

	var level zapcore.Level

	switch viper.GetString("logging.level") {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()

	var writes = []zapcore.WriteSyncer{write}

	if level == zap.DebugLevel {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	core := ecszap.NewCore(
		encoderConfig,
		zapcore.NewMultiWriteSyncer(writes...),
		level,
	)

	Logger = zap.New(
		core, zap.AddCaller(),
		zap.Fields(
			zap.String("application", viper.GetString("application.name")),
			zap.String("serverIp", util.LocalIp()),
			zap.Int("port", viper.GetInt("server.port")),
			zap.String("profile", viper.GetString("profiles.active")),
			zap.String("logTime", time.Now().Format(viper.GetString("date-time-format"))),
		),
	)
}
