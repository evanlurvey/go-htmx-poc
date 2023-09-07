package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

const sessionCookie = "session"

type sessionctxkey string

const sessionCtxKey sessionctxkey = "session"

type SessionState uint8

const (
	SessionState_Anonymous SessionState = iota
	SessionState_Authenticated
)

type Session struct {
	id      string
	token   string
	expires time.Time
	state   SessionState
	user    User
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) User() (User, bool) {
	if s.state == SessionState_Authenticated && time.Now().Before(s.expires) {
		return s.user, true
	}
	return User{}, false
}

func NewAnonymousSession() Session {
	return Session{
		id:    NewID(),
        token: NewID(32),
		state: SessionState_Anonymous,
	}
}

func NewAuthenticatedSession(u User) Session {
	return Session{
		id:    NewID(),
        token: NewID(32),
		state: SessionState_Authenticated,
		user:  u,
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
			Value:    session.token,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour),
			Secure:   true,
			HTTPOnly: true, // still sent with ajax calls
			SameSite: "Lax",
		})
	} else {
		// TODO: Validate and build session
		session = Session{
			id: sid,
		}
	}
	ctx = context.WithValue(ctx, sessionCtxKey, session)
	c.SetUserContext(ctx)
	return c.Next()
}
