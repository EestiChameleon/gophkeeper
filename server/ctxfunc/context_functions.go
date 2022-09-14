package ctxfunc

import (
	"context"
)

type ctxkey string

var (
	userID ctxkey = "userID"
)

// GetUserIDFromCTX returns from context userID if found.
func GetUserIDFromCTX(ctx context.Context) int {
	value, ok := ctx.Value(userID).(int)
	if !ok {
		return -1
	}
	return value
}

// SetUserIDToCTX add userID to the context.
func SetUserIDToCTX(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, userID, value)
}
