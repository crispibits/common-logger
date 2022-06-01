package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ENV_PROFILE string = "PROFILE"
const serverProfile string = "server"
const consoleProfile string = "console"

type Logger struct {
	*zap.SugaredLogger
}

// This config works well with GKE stackdriver logging
var serverConfig zap.Config = zap.Config{
	Encoding:         "json",
	Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "@timestamp",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	},
}

func New(opts ...zap.Option) *Logger {
	var cfg zap.Config
	profile := os.Getenv(ENV_PROFILE)
	switch profile {
	case consoleProfile:
		cfg = zap.NewDevelopmentConfig()
	case serverProfile:
		cfg = serverConfig
	default:
		cfg = serverConfig
	}

	logger, _ := cfg.Build(opts...)
	l := &Logger{SugaredLogger: logger.Sugar()}
	return l
}
