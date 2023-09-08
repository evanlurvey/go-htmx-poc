package forms

import (
	"htmx-poc/app/csrf"

	"github.com/gofiber/fiber/v2"
)

func NewService(csrf csrf.Service) Service {
	return Service{
		csrf: csrf,
	}
}

type Service struct {
	csrf csrf.Service
}

func (fv Service) Parse(c *fiber.Ctx, formData interface{ GetCSRFToken() string }) error {
	ctx := c.UserContext()
	if err := c.BodyParser(formData); err != nil {
		return err
	}
	if err := csrf.VerifyToken(ctx, formData.GetCSRFToken()); err != nil {
		return err
	}
	return nil
}
