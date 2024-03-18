package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Setup() {
	// Logger setup
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
