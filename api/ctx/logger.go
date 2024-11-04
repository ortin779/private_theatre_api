package ctx

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type LoggerCtxKey string

const LoggerKey LoggerCtxKey = "logger"

var ErrInvalidLogger = errors.New("invalid logger value")

func WithLogger(c context.Context, value *zap.Logger) context.Context {
	ctx := context.WithValue(c, LoggerKey, value)
	return ctx
}

func GetLogger(c context.Context) *zap.Logger {
	val, ok := c.Value(LoggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return val
}
