package contacts

import "fmt"

type Contact struct {
	ID    string `inputType:"hidden"`
	First string `validate:"gt=1"`
	Last  string `validate:"gt=1"`
	Phone string
	Email string `validate:"email"`
}

func (Contact) Fart() {
	fmt.Println("fart")
}
