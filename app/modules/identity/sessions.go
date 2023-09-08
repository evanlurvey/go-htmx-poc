package identity

import (
	"htmx-poc/utils"
	"time"
)

func NewAnonymousSession() Session {
	return Session{
		id:      utils.NewID(),
		token:   utils.NewID(32),
		expires: time.Now().Add(time.Hour * 12),
		state:   SessionState_Anonymous,
	}
}

func NewAuthenticatedSession(u User) Session {
	return Session{
		id:      utils.NewID(),
		token:   utils.NewID(32),
		expires: time.Now().Add(time.Hour * 12),
		state:   SessionState_Authenticated,
		user:    u,
	}
}
