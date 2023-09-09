package web

import (
	"context"
	"embed"
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/app/forms"
	"htmx-poc/app/logging"
	"htmx-poc/app/modules/contacts"
	"htmx-poc/app/modules/identity"
	"htmx-poc/app/template"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
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

	web.Use(compress.New())
	web.Use(logging.Middleware(l))
	web.Use(csrfSvc.Middleware)
	web.Use(identity.SessionMiddleware(identityDB))

	// FIX: remove in prod
	app.SetupAutoReloadWS(web)

	contacts.NewRouter(engine, contactsDB, form).Setup(web)

	identity.NewRouter(engine, identityDB, form).Setup(web)

	web.Get("/", func(c *fiber.Ctx) error {
		return engine.Render(c, "pages/index.html", map[string]any{})
	})

	web.Get("/counter", func(c *fiber.Ctx) error {
		count++
		return templateEngine.RenderComponent(c, "counter.html", template.M{
			"count": count,
		})
	})

	web.Get("/search", func(c *fiber.Ctx) error {
		all := []string{"evan", "eric", "jane"}
		out := []string{}
		q := c.Query("name")
		for _, i := range all {
			if strings.Contains(i, q) {
				out = append(out, i)
			}

		}
		return templateEngine.RenderComponent(c, "search.html", template.M{
			"results": out,
		})
	})

	return web
}

//go:embed views
var templatesFS embed.FS
var templateEngine = template.NewTemplateEngine(templatesFS, "")
var count = 0

func CounterComponent(ctx context.Context) any {
	return func() (template.HTML, error) {
		count++
		return templateEngine.RenderComponentHTML(ctx, "counter.html", template.M{
			"count": count,
		})
	}
}

func init() {
	template.RegisterComponent("counter", CounterComponent)
}
