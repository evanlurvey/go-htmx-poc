package contacts

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
	return r.templates.Render(c, "pages/single-form-page.html", map[string]any{
		"form": CreateContactForm,
	})
}

func (r *Router) NewPOST(c *fiber.Ctx) error {
	var formData ContactFormData
	if err := r.form.Parse(c, &formData); err != nil {
		return err
	}

	req := formData.Contact
	if ve := validation.ValidateStruct(req); ve != nil {
		form := CreateContactForm.GenerateFields(req, ve)
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
		return r.templates.Render(c, "pages/not-found.html", nil)
	}
	return r.templates.Render(c, "pages/contacts/contacts-detail.html", map[string]any{
		"contact": contact,
	})
}

func (r *Router) EditGET(c *fiber.Ctx) error {
	id := c.Params("id")
	contact, found := r.db.GetContactByID(id)
	if !found {
		return r.templates.Render(c, "pages/not-found.html", nil)
	}
	form := UpdateContactForm.GenerateFields(contact)
	return r.templates.Render(c, "pages/single-form-page.html", map[string]any{
		"form": form,
	})
}

func (r *Router) EditPOST(c *fiber.Ctx) error {
	id := c.Params("id")
	var formData ContactFormData
	if err := r.form.Parse(c, &formData); err != nil {
		return err
	}
	req := formData.Contact
	if id != req.ID {
		return c.SendStatus(401)
	}
	if ve := validation.ValidateStruct(req); ve != nil {
		form := UpdateContactForm.GenerateFields(req, ve)
		return r.templates.Render(c, "pages/single-form-page.html", fiber.Map{
			"form": form,
		})
	}
	r.db.UpdateContactByID(req)
	return c.Redirect("/contacts", 303)
}
