package main

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/app/modules/auth"
	"htmx-poc/app/modules/contacts"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// TODO: current thought is I need to make a form input field struct that can be passed in
// to auto build forms and handle error responses and what not in a standard way.
// should be helpful to keep things consistent and quick.

var appversion string

func init() {
	appversion = app.NewID()
}

type controller interface {
	Setup(fiber.Router)
}

func setupController(app fiber.Router, c controller) {
	c.Setup(app)
}

func main() {
	web := fiber.New(fiber.Config{
		Immutable:             true, // string buffers get reused otherwise and shit gets weird when using in memory things
		DisableStartupMessage: true,
	})

	csrf := csrf.New([]byte("TODO: Change me"))

	engine := app.NewTemplateEngine(csrf, appversion, "layouts/main.html")

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

	contactsDB := contacts.NewDB()
	contacts.NewRouter(engine, contactsDB, csrf).Setup(web)

	authDB := auth.NewDB()
	auth.NewRouter(engine, authDB, csrf).Setup(web)

	slog.Info("starting server")
	if err := web.Listen(":8080"); err != nil {
		slog.Error("failed to start server", err)
	}
}
