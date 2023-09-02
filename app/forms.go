package app

import (
	"context"
	"htmx-poc/app/csrf"
	"htmx-poc/validation"
	"reflect"
	"slices"

	"github.com/gofiber/fiber/v2"
)

type FormField struct {
	Label string
	Name  string
	Value any
	Type  string
	Error string
	// gonna have to add other attrs n stuff
	// start of a component tho
}

type FormFields []FormField

type Form struct {
	Title            string
	Template         string
	Fields           FormFields
	SubmitButtonText string
	BackButton       bool
	CSRFToken        string
}

func (f Form) Clone() Form {
	f.Fields = slices.Clone(f.Fields)
	return f
}

func (f Form) AddCSRFToken(ctx context.Context, csrf csrf.Service) Form {
	session := SessionFromCtx(ctx)
	f.CSRFToken = csrf.NewToken(session.ID)
	return f
}

// returns a new cloned form with fields including prev data and errors
func (f Form) GenerateFields(in any, ve ...validation.Errors) Form {
	form := f.Clone()
	v := reflect.ValueOf(in)
	t := v.Type()
	fc := v.NumField()
	form.Fields = make(FormFields, fc)
	for i := 0; i < fc; i++ {
		f := v.Field(i)
		ft := t.Field(i)
		var (
			label        = OrDefault(ft.Tag.Get("label"), ft.Name)
			name         = OrDefault(ft.Tag.Get("name"), ft.Name)
			value        any
			inputType    = OrDefault(ft.Tag.Get("inputType"), "text")
			errorMessage string
		)
		// don't sendback sensitive data
		if inputType != "password" || ft.Tag.Get("sensitive") == "true" {
			value = f.Interface()
		}
		// check validation errors
		if len(ve) == 1 {
			for _, ve := range ve[0] {
				if ve.Name == name {
					errorMessage = ve.Message
				}
			}
		}
		// build
		form.Fields[i] = FormField{
			Label: label,
			Name:  name,
			Value: value,
			Type:  inputType,
			Error: errorMessage,
		}
	}
	return form
}

func OrDefault[T comparable](v, d T) T {
	var empty T
	if v == empty {
		return d
	}
	return v
}

func NewFormService(csrf csrf.Service) FormService {
	return FormService{
		csrf: csrf,
	}
}

type FormService struct {
	csrf csrf.Service
}

func (fv FormService) Parse(c *fiber.Ctx, formData interface{ GetCSRFToken() string }) error {
	ctx := c.UserContext()
	session := SessionFromCtx(ctx)
	if err := c.BodyParser(formData); err != nil {
		return err
	}
	if err := fv.csrf.VerifyToken(session.ID, formData.GetCSRFToken()); err != nil {
		return err
	}
	return nil
}
