package auth

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(templates app.TemplateEngine, db *DB, csrf csrf.Service) *Router {
	return &Router{
		templates: templates,
		db:        db,
		csrf:      csrf,
	}
}

type Router struct {
	templates app.TemplateEngine
	db        *DB
	csrf      csrf.Service
}

func (r *Router) Setup(rtr fiber.Router) {
	rtr = rtr.Group("/accounts")
	rtr.Get("/", r.LoginGET)
}

func (r *Router) LoginGET(c *fiber.Ctx) error {
	ctx := c.UserContext()
	form := LoginForm.AddCSRFToken(ctx, r.csrf)
	return r.templates.Render(c, "pages/login.html", map[string]any{
		"form": form,
	})
}
