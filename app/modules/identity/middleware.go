package identity

import (
	"htmx-poc/app/logging"
	"htmx-poc/app/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const sessionCookie = "session"

// ensures there is a session associated with the request
func SessionMiddleware(db *DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		token := c.Cookies(sessionCookie)

		// validate session
		session, err := db.getSessionByToken(ctx, token)
		if err != nil {
			// TODO: Log error or something
		}
		user, _ := session.User()

		// fallback to anonymous session
		if session.id == "" || !session.Valid() {
			session = newAnonymousSession()
			db.storeSession(ctx, session)
			c.Cookie(&fiber.Cookie{
				Name:     sessionCookie,
				Value:    session.token,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour),
				Secure:   true,
				HTTPOnly: true, // still sent with ajax calls
				SameSite: "Lax",
			})
		}
		ctx = SessionWithContext(ctx, session)
		ctx = logging.QWithContext(ctx, zap.String("session_id", session.id))
		// attach in global template context
		ctx = template.WithContext(ctx, template.M{
			"_user": user,
		})
		c.SetUserContext(ctx)
		return c.Next()
	}
}
