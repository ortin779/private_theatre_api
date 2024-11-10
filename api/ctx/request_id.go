package ctx

import (
	"context"

	"github.com/google/uuid"
)

type RequestIdCtxKey string

const RequestIdKey RequestIdCtxKey = "req-id"

func WithRequestId(c context.Context) context.Context {
	ctx := context.WithValue(c, RequestIdKey, uuid.NewString())
	return ctx
}

func GetRequestId(c context.Context) string {
	val, ok := c.Value(RequestIdKey).(string)
	if !ok {
		return uuid.NewString()
	}
	return val
}
