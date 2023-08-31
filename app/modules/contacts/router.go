package contacts

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/validation"

	"github.com/gofiber/fiber/v2"
)

// TODO: I really don't like this or how I am parsing and validating forms 
// find a better way to do this
func parsemedaddy(c *fiber.Ctx, formData interface{ GetCSRFToken() string }, verifier interface{ VerifyToken(string, string) error }) error {
	if err := c.BodyParser(formData); err != nil {
		return err
	}
	if err := verifier.VerifyToken(app.SessionFromCtx(c.UserContext()).ID, formData.GetCSRFToken()); err != nil {
		return err
	}
	return nil
}

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
	rtr = rtr.Group("/contacts")
	rtr.Get("/", r.Index)
	rtr.Get("/new", r.NewGET)
	rtr.Post("/new", r.NewPOST)
	rtr.Get("/:id", r.Details)
	rtr.Get("/:id/edit", r.EditGET)
	rtr.Post("/:id/edit", r.EditPOST)

}

func (r *Router) Index(c *fiber.Ctx) error {
	return r.templates.Render(c, "pages/contacts/contacts.html", map[string]any{
		"contacts": r.db.FindContacts(),
	})
}

func (r *Router) NewGET(c *fiber.Ctx) error {
	ctx := c.UserContext()
	form := CreateContactForm.AddCSRFToken(ctx, r.csrf)
	return r.templates.Render(c, "pages/single-form-page.html", map[string]any{
		"form": form,
	})
}

func (r *Router) NewPOST(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var formData ContactFormData
	if err := parsemedaddy(c, &formData, r.csrf); err != nil {
		return err
	}

	req := formData.Contact
	if ve := validation.ValidateStruct(req); ve != nil {
		form := CreateContactForm.AddCSRFToken(ctx, r.csrf).GenerateFields(req, ve)
		form.Title = "get it right"
		return r.templates.Render(c, "pages/single-form-page.html", fiber.Map{
			"form": form,
		})
	}
	// save contact
	req.ID = app.NewID()
	r.db.AddContact(req)
	return c.Redirect("/contacts", 303)
}

func (r *Router) Details(c *fiber.Ctx) error {
	id := c.Params("id")
	contact, found := r.db.GetContactByID(id)
	if !found {
		return c.SendStatus(404)
	}
	return r.templates.Render(c, "pages/contacts/contacts-detail.html", map[string]any{
		"contact": contact,
	})
}

func (r *Router) EditGET(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	contact, found := r.db.GetContactByID(id)
	if !found {
		return c.SendStatus(404)
	}
	form := UpdateContactForm.AddCSRFToken(ctx, r.csrf).GenerateFields(contact)
	return r.templates.Render(c, "pages/single-form-page.html", map[string]any{
		"id":   id,
		"form": form,
	})
}

func (r *Router) EditPOST(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	var reqBody ContactFormData
	if err := parsemedaddy(c, &reqBody, r.csrf); err != nil {
		return err
	}
	req := reqBody.Contact
	if id != req.ID {
		return c.SendStatus(401)
	}
	if ve := validation.ValidateStruct(req); ve != nil {
		form := UpdateContactForm.AddCSRFToken(ctx, r.csrf).GenerateFields(req, ve)
		return r.templates.Render(c, "pages/single-form-page.html", fiber.Map{
			"form": form,
		})
	}
	r.db.UpdateContactByID(req)
	return c.Redirect("/contacts", 303)
}
