package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func New() *zap.SugaredLogger {
	if logger == nil {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		logger = zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderCfg),
				zapcore.AddSync(os.Stdout),
				zap.NewAtomicLevelAt(zap.DebugLevel),
			),
		).Sugar()
	}

	return logger
}
