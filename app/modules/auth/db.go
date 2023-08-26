package auth

import "sync"

type DB struct {
	m sync.RWMutex
}

func NewDB() *DB {
	return &DB{}
}
