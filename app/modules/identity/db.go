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
	m            sync.RWMutex
	users        []user
	loginAttempt []loginAttempt
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) storeUser(u user) error {
	db.m.Lock()
	defer db.m.Unlock()

	idx := slices.IndexFunc(db.users, func(uu user) bool {
		return u.id == uu.id
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
		return id == u.id
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
		return email == u.email
	})
	if idx < 0 {
		return user{}, nil
	}
	return db.users[idx], nil
}

func (db *DB) createLoginAttempt(ctx context.Context, la loginAttempt) error {
	db.m.Lock()
	defer db.m.Unlock()

	if slices.ContainsFunc(db.loginAttempt, func(l loginAttempt) bool { return la.id == l.id }) {
		return fmt.Errorf("id conflict")
	}
	db.loginAttempt = append(db.loginAttempt, la)
	return nil
}

func (db *DB) getRecentAttempts(ctx context.Context, email string, last time.Duration) (loginAttempts, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	window := time.Now().Add(-last)
	attempts := utils.FilterFunc(db.loginAttempt, func(la loginAttempt) bool {
		return la.at.After(window) && la.email == email
	})
	return attempts, nil
}
