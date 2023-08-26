package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

const sessionCookie = "session"

type sessionctxkey string

const sessionCtxKey sessionctxkey = "session"

type SessionType string

const (
	SessionTypeAnonymous = "anonymous"
)

type Session struct {
	ID   string
	Type string
}

func NewAnonymousSession() Session {
	return Session{
		ID:   NewID(),
		Type: SessionTypeAnonymous,
	}
}

func SessionFromCtx(ctx context.Context) Session {
	session, ok := ctx.Value(sessionCtxKey).(Session)
	if !ok {
		return NewAnonymousSession()
	}
	return session
}

// ensures there is a session associated with the request
func SessionMiddleware(c *fiber.Ctx) error {
	ctx := c.UserContext()
	sid := c.Cookies(sessionCookie)
	var session Session
	if sid == "" {
		session = NewAnonymousSession()
		c.Cookie(&fiber.Cookie{
			Name:     sessionCookie,
			Value:    session.ID,
			Path:     "/",
			MaxAge:   int(time.Hour * 24),
			Secure:   true,
			HTTPOnly: true, // still sent with ajax calls
			SameSite: "Lax",
		})
	} else {
		// TODO: Validate and build session
		session = Session{
			ID:   sid,
			Type: SessionTypeAnonymous,
		}
	}
	ctx = context.WithValue(ctx, sessionCtxKey, session)
	c.SetUserContext(ctx)
	return c.Next()
}
