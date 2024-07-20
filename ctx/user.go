package ctx

import (
	"context"
	"errors"
)

type UserIdKey string

var UserIdCtxKey UserIdKey = "userId"

var (
	ErrInvalidUserId = errors.New("invalid user id type")
)

func WithUserId(c context.Context, value string) context.Context {
	ctx := context.WithValue(c, UserIdCtxKey, value)
	return ctx
}

func UserIdValue(c context.Context) (string, error) {
	val, ok := c.Value(UserIdCtxKey).(string)
	if !ok {
		return "", ErrInvalidUserId
	}
	return val, nil
}
