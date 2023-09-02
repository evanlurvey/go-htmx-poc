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
	Email    string `validate:"email"`
	Password string `validate:"required" inputType:"password"`
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
	FirstName       string `validate:"required" label:"First Name"`
	LastName        string `validate:"required" label:"Last Name"`
	Email           string `validate:"email"`
	Password        string `validate:"gte=8,lte=70" inputType:"password"`
	ConfirmPassword string `validate:"eqfield=Password" label:"Confirm Password" inputType:"password"`
}

func init() {
	CreateAccountForm = CreateAccountForm.GenerateFields(CreateAccountRequest{})
}
