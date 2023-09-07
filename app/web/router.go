package web

import (
	"context"
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/app/modules/contacts"
	"htmx-poc/app/modules/identity"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
    // TODO: Have and use a config lol
    CSRFSecret []byte
}

func NewWebApp(ctx context.Context, cfg Config) *fiber.App {
    l := app.Logger(ctx)
	web := fiber.New(fiber.Config{
		Immutable:             true, // string buffers get reused otherwise and shit gets weird when using in memory things
		DisableStartupMessage: true,
	})

	csrf := csrf.New(cfg.CSRFSecret)

	engine := app.NewTemplateEngine(csrf, "layouts/main.html")

	web.Use(app.LoggingMiddleware(l))
	web.Use(csrf.ErrorHandler)
	web.Use(app.SessionMiddleware)

	// FIX: remove in prod
	app.SetupAutoReloadWS(web)

	web.Get("/", func(c *fiber.Ctx) error {
		return engine.Render(c, "pages/index.html", map[string]any{
			"ctx":  c.UserContext(),
			"name": "evan<p>lol</p>",
		})
	})

	form := app.NewFormService(csrf)

	contactsDB := contacts.NewDB()
	contacts.NewRouter(engine, contactsDB, form).Setup(web)

	authDB := identity.NewDB()
	identity.NewRouter(engine, authDB, form).Setup(web)

    return web
}
