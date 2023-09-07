package main

import (
	"context"
	"htmx-poc/app"
	"htmx-poc/app/web"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	l, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	return l
}

func main() {
	l := logger()
	ctx := app.AttachLogger(context.Background(), l)

	web := web.NewWebApp(ctx, web.Config{
		CSRFSecret: []byte("not a secret"),
	})

	l.Info("starting server")
	if err := web.Listen(":8080"); err != nil {
		l.Panic("failed to start server", zap.Error(err))
	}
}
