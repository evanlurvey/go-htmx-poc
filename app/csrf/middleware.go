package csrf

import (
	"errors"
	"htmx-poc/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

const csrfSessionCookie = "csst"

func (s Service) Middleware(c *fiber.Ctx) error {
	ctx := c.UserContext()
	csst := c.Cookies(csrfSessionCookie)
	// generate and save cookie if we don't have one
	if csst == "" || len(csst) > 500 {
		csst = utils.NewID(32)
		c.Cookie(&fiber.Cookie{
			Name:     csrfSessionCookie,
			Value:    csst,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   true,
			HTTPOnly: true, // still sent with ajax calls
			SameSite: "Lax",
		})
	}
	ctx = WithContext(ctx, s.NewToken(csst))
	c.SetUserContext(ctx)
	// run requests with csrf
	err := c.Next()
	// check errors
	if errors.Is(err, InvalidCSRFError) {
		// TODO: Log these errors
		return c.Status(401).SendString("invalid csrf")
	}
	return err
}
