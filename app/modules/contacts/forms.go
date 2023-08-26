package contacts

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
)

type ContactFormData struct {
	csrf.FormData
	Contact
}

var (
	CreateContactForm = app.Form{
		Title:            "Create Contact",
		SubmitButtonText: "Create",
		BackButton:       true,
	}
	UpdateContactForm = app.Form{
		Title:            "Update Contact",
		SubmitButtonText: "Update",
		BackButton:       true,
	}
)

func init() {
	CreateContactForm = CreateContactForm.GenerateFields(Contact{})
	UpdateContactForm = UpdateContactForm.GenerateFields(Contact{})
}
