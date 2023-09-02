package main

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/app/modules/auth"
	"htmx-poc/app/modules/contacts"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: current thought is I need to make a form input field struct that can be passed in
// to auto build forms and handle error responses and what not in a standard way.
// should be helpful to keep things consistent and quick.

var appversion string

func init() {
	appversion = app.NewID()
}

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

	web := fiber.New(fiber.Config{
		Immutable:             true, // string buffers get reused otherwise and shit gets weird when using in memory things
		DisableStartupMessage: true,
	})

	csrf := csrf.New([]byte("TODO: Change me"))

	engine := app.NewTemplateEngine(csrf, appversion, "layouts/main.html")

	web.Use(app.LoggingMiddleware(l))
	web.Use(csrf.ErrorHandler)
	web.Use(app.SessionMiddleware)

	// FIX: remove in prod
	app.SetupAutoReloadWS(web, appversion)

	web.Get("/", func(c *fiber.Ctx) error {
		return engine.Render(c, "pages/index.html", map[string]any{
			"ctx":  c.UserContext(),
			"name": "evan<p>lol</p>",
		})
	})

	form := app.NewFormService(csrf)

	contactsDB := contacts.NewDB()
	contacts.NewRouter(engine, contactsDB, form).Setup(web)

	authDB := auth.NewDB()
	auth.NewRouter(engine, authDB, form).Setup(web)

	l.Info("starting server")
	if err := web.Listen(":8080"); err != nil {
		l.Panic("failed to start server", zap.Error(err))
	}
}
