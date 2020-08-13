package b_globals

import (
	"context"
)

type requestIdKey string

var (
	reqID = requestIdKey("req-id")
)

func GetRequestID(ctx context.Context) (string, bool) {
	id, exists := ctx.Value(reqID).(string)
	return id, exists
}

func WithRequestID(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, reqID, reqId)
}
