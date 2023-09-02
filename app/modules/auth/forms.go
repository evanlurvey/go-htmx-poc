package auth

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
)

// Login

var LoginForm = app.Form{
	Title:            "Login",
	SubmitButtonText: "Login",
	BackButton:       true,
}

func init() {
	LoginForm = LoginForm.GenerateFields(LoginRequest{})
}

type LoginFormData struct {
	csrf.FormData
	LoginRequest
}

type LoginRequest struct {
	Email    string
	Password string `inputType:"password"`
}

// Create Account

var CreateAccountForm = app.Form{
	Title:            "Create Account",
	SubmitButtonText: "Create Account",
	BackButton:       true,
}

type CreateAccountFormData struct {
	csrf.FormData
	CreateAccountRequest
}

type CreateAccountRequest struct {
	Email    string
	Password string `inputType:"password"`
}
