package web

import (
	"context"
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/app/forms"
	"htmx-poc/app/logging"
	"htmx-poc/app/modules/contacts"
	"htmx-poc/app/modules/identity"
	"htmx-poc/app/template"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	// TODO: Have and use a config lol
	CSRFSecret []byte
}

func NewWebApp(ctx context.Context, cfg Config) *fiber.App {
	l := logging.FromContext(ctx)
	web := fiber.New(fiber.Config{
		Immutable:             true, // string buffers get reused otherwise and shit gets weird when using in memory things
		DisableStartupMessage: true,
	})

	csrfSvc := csrf.New(cfg.CSRFSecret)

	// DB Setup
	contactsDB := contacts.NewDB()
	identityDB := identity.NewDB()
	// Service Setup
	engine := template.NewTemplateEngine(nil, "layouts/main.html")
	form := forms.NewService(csrfSvc)

	web.Use(logging.Middleware(l))
	web.Use(csrfSvc.Middleware)
	web.Use(identity.SessionMiddleware(identityDB))

	// FIX: remove in prod
	app.SetupAutoReloadWS(web)

	contacts.NewRouter(engine, contactsDB, form).Setup(web)

	identity.NewRouter(engine, identityDB, form).Setup(web)

	count := 0
	web.Get("/counter", func(c *fiber.Ctx) error {
		count++
		return engine.Render(c, "pages/counter.html", template.M{
			"count": count,
		})
	})

	return web
}
