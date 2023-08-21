package main

type Contact struct {
	ID    string
	First string `validate:"gt=1"`
	Last  string `validate:"gt=1"`
	Phone string
	Email string `validate:"email"`
}
