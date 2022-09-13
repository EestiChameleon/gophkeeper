package ctxfunc

import (
	"context"
)

type ctxkey string

var (
	userID ctxkey = "userID"
	uID    int
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

// SetUserID saves userID to local variable. Used for login handler.
func SetUserID(value int) {
	uID = value
}

// GetUserID returns userID from local variable. Used for login handler.
func GetUserID() int {
	return uID
}
