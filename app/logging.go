package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type loggingctxkey string

const loggingCtxKey loggingctxkey = "logging"

func AttachLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggingCtxKey, l)
}

func Logger(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggingCtxKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return l
}

func LoggingMiddleware(l *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		traceID := NewID(4)
		start := time.Now()

		// attach new logger to context
		l := l.With(zap.String("trace_id", traceID))
		ctx = AttachLogger(ctx, l)
		c.SetUserContext(ctx)

		// chain
		l.Debug(
			"request starting",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)
		err := c.Next()
		l.Info(
			"request finished",
            // req fields should match above
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
            // resp fields
            zap.Int("status_code", c.Response().StatusCode()),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
		return err
	}
}
