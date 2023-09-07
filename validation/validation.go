package validation

import (
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate            *validator.Validate = validator.New()
	validate_translator ut.Translator
)

func init() {
	en := en.New()
	uni := ut.New(en, en)

	// Ignore if found or not.
	validate_translator, _ = uni.GetTranslator("en")

	// Get proper field names from json
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// the split is for when there is extra shit like omitempty
		name := fld.Name
		if customName := fld.Tag.Get("label"); customName != "" {
			name = customName
		}
		return name
	})

	// Register translations for pretty text
	if err := en_translations.RegisterDefaultTranslations(validate, validate_translator); err != nil {
		panic(err)
	}
}

type Errors []Error

type Error struct {
	Name    string
	Message string
}

// Returns a ValidationErrors error type.
func ValidateStruct(data any) Errors {
	err := validate.Struct(data)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		out := make(Errors, len(errs))
		for _, ve := range errs {
			out = append(out, Error{
				Name:    ve.StructField(),
				Message: ve.Translate(validate_translator),
			})
		}
		return out
	}
	return nil
}
