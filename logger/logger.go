package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       os.Getenv("APP_ENV") != "development",
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling:          nil,
		EncoderConfig:     encoderCfg,
		Encoding:          "json",
		OutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getegid(),
		},
	}

	return zap.Must(cfg.Build())
}
