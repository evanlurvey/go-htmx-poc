package main

import (
	"htmx-poc/validation"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type ContactsRouter struct {
	templates *TemplateEngine
	db        *DB
}

func (r *ContactsRouter) Setup(rtr fiber.Router) {
	rtr = rtr.Group("/contacts")
	rtr.Get("/", r.Index)
	rtr.Get("/new", r.NewGET)
	rtr.Post("/new", r.NewPOST)
	rtr.Get("/:id", r.Details)
	rtr.Get("/:id/edit", r.EditGET)
	rtr.Post("/:id/edit", r.EditPOST)

}

func (r *ContactsRouter) Index(c *fiber.Ctx) error {
	return r.templates.Render(c, "contacts.html", map[string]any{
		"contacts": r.db.FindContacts(),
	})
}

type Input struct {
	Label string
	Name  string
	Value any
	Type  string
	Error string
	// gonna have to add other attrs n stuff
	// start of a component tho
}

type FormInputs map[string]Input

func (f FormInputs) SetInput(in Input) {
	f[in.Name] = in
}

func (f FormInputs) AddTextInput(label, name, value string) {
	f.SetInput(Input{
		Label: label,
		Name:  name,
		Value: value,
		Type:  "text",
		Error: "",
	})
}
func (f FormInputs) SetError(name, error string) {
	if i, found := f[name]; found {
		i.Error = error
		f[name] = i
	}
}
func (f FormInputs) SetValidationErrors(ve validation.Errors) {
	for _, e := range ve {
		f.SetError(e.Name, e.Message)
	}
}

func (r *ContactsRouter) NewGET(c *fiber.Ctx) error {
	form := FormInputs{}
	form.AddTextInput("First", "First", "")
	form.AddTextInput("Last", "Last", "")
	form.AddTextInput("Phone", "Phone", "")
	form.AddTextInput("Email", "Email", "")
	return r.templates.Render(c, "contacts-form.html", map[string]any{
		"form": form,
		"new":  true,
	})
}

func (r *ContactsRouter) NewPOST(c *fiber.Ctx) error {
	var reqBody struct {
		CSRFToken string `form:"csrf_token"`
		Contact
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}
	if !VerifyCSRFToken("sessionid", reqBody.CSRFToken) {
		return c.Status(401).SendString("invalid csrf")
	}
	req := reqBody.Contact
	if ve := validation.ValidateStruct(req); ve != nil {
		form := FormInputs{}
		form.AddTextInput("First", "First", req.First)
		form.AddTextInput("Last", "Last", req.Last)
		form.AddTextInput("Phone", "Phone", req.Phone)
		form.AddTextInput("Email", "Email", req.Email)
		form.SetValidationErrors(ve)
		return r.templates.Render(c, "contacts-form.html", fiber.Map{
			"form": form,
			"new":  true,
		})
	}
	req.ID = newID()
	r.db.AddContact(req)
	return c.Redirect("/contacts", 303)
}

func (r *ContactsRouter) Details(c *fiber.Ctx) error {
	id := c.Params("id")
	slog.Info("contact detail", slog.String("contact_id", id))
	contact, found := r.db.GetContactByID(id)
	if !found {
		return c.SendStatus(404)
	}
	return r.templates.Render(c, "contacts-detail.html", map[string]any{
		"contact": contact,
	})
}

func (r *ContactsRouter) EditGET(c *fiber.Ctx) error {
	id := c.Params("id")
	slog.Info("contact detail", slog.String("contact_id", id))
	contact, found := r.db.GetContactByID(id)
	if !found {
		return c.SendStatus(404)
	}
	form := FormInputs{}
	form.AddTextInput("First", "First", contact.First)
	form.AddTextInput("Last", "Last", contact.Last)
	form.AddTextInput("Phone", "Phone", contact.Phone)
	form.AddTextInput("Email", "Email", contact.Email)
	return r.templates.Render(c, "contacts-form.html", map[string]any{
		"id":   id,
		"form": form,
	})
}
func (r *ContactsRouter) EditPOST(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		CSRFToken string `form:"csrf_token"`
		Contact
	}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if !VerifyCSRFToken("sessionid", req.CSRFToken) {
		return c.Status(401).SendString("invalid csrf")
	}
    if id != req.ID {
        return c.SendStatus(401)
    }
	if ve := validation.ValidateStruct(req.Contact); ve != nil {
		form := FormInputs{}
		form.AddTextInput("First", "First", req.First)
		form.AddTextInput("Last", "Last", req.Last)
		form.AddTextInput("Phone", "Phone", req.Phone)
		form.AddTextInput("Email", "Email", req.Email)
		form.SetValidationErrors(ve)
		return r.templates.Render(c, "contacts-form.html", fiber.Map{
			"id":   id,
			"form": form,
		})
	}
	r.db.UpdateContactByID(req.Contact)
	return c.Redirect("/contacts", 303)
}
