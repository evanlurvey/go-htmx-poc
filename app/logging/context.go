package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggingctxkey string

const loggingCtxKey loggingctxkey = "logging"

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggingCtxKey, l)
}

func QWithContext(ctx context.Context, fields ...zap.Field) context.Context {
	return WithContext(ctx, FromContext(ctx).With(fields...))
}

func FromContext(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggingCtxKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return l
}
