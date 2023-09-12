package main

import (
	"context"
	"htmx-poc/app/logging"
	"htmx-poc/app/web"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate tailwindcss -o app/web/static/style.css -m
// TODO: use go generate to scrape together all of the localization entries and build json for it

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
	csrfSecret := os.Getenv("SECRET_CSRF_TOKEN_KEY")
	appEnv := os.Getenv("APP_ENV")
	l := logger().With(zap.String("app_env", appEnv))
	l.Info("csrfDebug", zap.String("csrf", csrfSecret))
	ctx := logging.WithContext(context.Background(), l)

	web := web.NewWebApp(ctx, web.Config{
		CSRFSecret: []byte(csrfSecret),
	})

	l.Info("starting server")
	if err := web.Listen(":8080"); err != nil {
		l.Panic("failed to start server", zap.Error(err))
	}
}
