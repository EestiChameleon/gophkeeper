package ctxfunc

import (
	"context"
	"errors"
)

type ctxkey string

var (
	userID ctxkey = "userID"
	uID    int
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
)

func GetUserIDFromCTX(ctx context.Context) int {
	value, ok := ctx.Value(userID).(int)
	if !ok {
		return -1
	}
	return value
}

func SetUserIDToCTX(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, userID, value)
}

func SetUserID(value int) {
	uID = value
}

func GetUserID() int {
	return uID
}
