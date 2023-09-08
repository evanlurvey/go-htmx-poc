package logging

import (
	"htmx-poc/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Middleware(l *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		traceID := utils.NewID(4)
		start := time.Now()

		// attach new logger to context
		l := l.With(zap.String("trace_id", traceID))
		ctx = WithContext(ctx, l)
		c.SetUserContext(ctx)

		// chain
		l.Debug(
			"request starting",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)
		err := c.Next()

		// after
		// extract logger, might have new goodies
		l = FromContext(c.UserContext())
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
