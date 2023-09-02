package auth

import (
	"htmx-poc/app"
	"htmx-poc/validation"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(templates app.TemplateEngine, db *DB, form app.FormService) *Router {
	return &Router{
		templates: templates,
		db:        db,
		form:      form,
	}
}

type Router struct {
	templates app.TemplateEngine
	db        *DB
	form      app.FormService
}

func (r *Router) Setup(rtr fiber.Router) {
	rtr = rtr.Group("/identity")
	rtr.Get("/login", r.LoginGET)
	rtr.Post("/login", r.LoginPOST)
	rtr.Get("/create-account", r.CreateAccountGET)
	rtr.Post("/create-account", r.CreateAccountPOST)
}

func (r *Router) LoginGET(c *fiber.Ctx) error {
	return r.templates.Render(c, "pages/login.html", map[string]any{
		"form": LoginForm,
	})
}

func (r *Router) LoginPOST(c *fiber.Ctx) error {
	var formData LoginFormData
	if err := r.form.Parse(c, &formData); err != nil {
		return err
	}

	req := formData.LoginRequest
	if ve := validation.ValidateStruct(req); ve != nil {
		form := LoginForm.GenerateFields(req, ve)
		return r.templates.Render(c, "pages/login.html", fiber.Map{
			"form": form,
		})
	}
	return c.Redirect("/", 303)
}

func (r *Router) CreateAccountGET(c *fiber.Ctx) error {
	return r.templates.Render(c, "pages/login.html", map[string]any{
		"form": CreateAccountForm,
	})
}

func (r *Router) CreateAccountPOST(c *fiber.Ctx) error {
	var formData CreateAccountFormData
	if err := r.form.Parse(c, &formData); err != nil {
		return err
	}

	req := formData.CreateAccountRequest
	if ve := validation.ValidateStruct(req); ve != nil {
		form := CreateAccountForm.GenerateFields(req, ve)
		return r.templates.Render(c, "pages/login.html", fiber.Map{
			"form": form,
		})
	}

	return c.Redirect("/identity/login", 303)
}
