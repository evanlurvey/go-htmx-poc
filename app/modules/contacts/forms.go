package contacts

import (
	"htmx-poc/app/csrf"
	"htmx-poc/app/forms"
)

type ContactFormData struct {
	csrf.FormData
	Contact
}

var (
	CreateContactForm = forms.Form{
		Title:            "Create Contact",
		SubmitButtonText: "Create",
		BackButton:       true,
	}
	UpdateContactForm = forms.Form{
		Title:            "Update Contact",
		SubmitButtonText: "Update",
		BackButton:       true,
	}
)

func init() {
	CreateContactForm = CreateContactForm.GenerateFields(Contact{})
	UpdateContactForm = UpdateContactForm.GenerateFields(Contact{})
}
