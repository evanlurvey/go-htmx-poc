package auth

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
)

var (
	LoginForm = app.Form{
		Title:            "Login",
		SubmitButtonText: "Login",
		BackButton:       true,
	}
)

type LoginFormData struct {
	csrf.FormData
	LoginRequest
}

type LoginRequest struct {
	Email    string
	Password string `inputType:"password"`
}

func init() {
	LoginForm = LoginForm.GenerateFields(LoginRequest{})
}
