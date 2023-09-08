package identity

import (
	"context"
	"fmt"
	"htmx-poc/utils"
	"slices"
	"sync"
	"time"
)

type DB struct {
	m             sync.RWMutex
	users         []user
	loginAttempts []loginAttempt
	sessions      []Session
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) storeUser(u user) error {
	db.m.Lock()
	defer db.m.Unlock()

	idx := slices.IndexFunc(db.users, func(uu user) bool {
		return u.ID == uu.ID
	})
	if idx >= 0 {
		db.users[idx] = u
	} else {
		db.users = append(db.users, u)
	}
	return nil
}

func (db *DB) getUserByID(id string) (user, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	idx := slices.IndexFunc(db.users, func(u user) bool {
		return id == u.ID
	})
	if idx < 0 {
		return user{}, nil
	}
	return db.users[idx], nil
}

func (db *DB) getUserByEmail(ctx context.Context, email string) (user, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	idx := slices.IndexFunc(db.users, func(u user) bool {
		return email == u.Email
	})
	if idx < 0 {
		return user{}, nil
	}
	return db.users[idx], nil
}

func (db *DB) createLoginAttempt(ctx context.Context, la loginAttempt) error {
	db.m.Lock()
	defer db.m.Unlock()

	if slices.ContainsFunc(db.loginAttempts, func(l loginAttempt) bool { return la.id == l.id }) {
		return fmt.Errorf("id conflict")
	}
	db.loginAttempts = append(db.loginAttempts, la)
	return nil
}

func (db *DB) getRecentAttempts(ctx context.Context, email string, last time.Duration) (loginAttempts, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	window := time.Now().Add(-last)
	attempts := utils.FilterFunc(db.loginAttempts, func(la loginAttempt) bool {
		return la.at.After(window) && la.email == email
	})
	return attempts, nil
}

func (db *DB) storeSession(ctx context.Context, s Session) error {
	db.m.Lock()
	defer db.m.Unlock()

	idx := slices.IndexFunc(db.sessions, func(ss Session) bool {
		return s.id == ss.id
	})
	if idx >= 0 {
		db.sessions[idx] = s
	} else {
		db.sessions = append(db.sessions, s)
	}
	return nil
}

func (db *DB) getSessionByToken(ctx context.Context, token string) (Session, error) {
	db.m.RLock()
	defer db.m.RUnlock()
	idx := slices.IndexFunc(db.sessions, func(s Session) bool {
		return token == s.token
	})
	if idx < 0 {
		return Session{}, nil
	}
	session := db.sessions[idx]

	if session.user.ID == "" {
		return session, nil
	}

	idx = slices.IndexFunc(db.users, func(u user) bool {
		return session.user.ID == u.ID
	})
	if idx < 0 {
		// session with a user set that doesn't exist. shouldn't ever happen
		return Session{}, nil
	}
	session.user = db.users[idx].User
	return session, nil

}
