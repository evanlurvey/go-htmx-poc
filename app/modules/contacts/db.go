package contacts

import (
	"sync"
)

type DB struct {
	m        sync.RWMutex
	contacts []Contact
}

func NewDB() *DB {
	db := &DB{
		contacts: []Contact{
			{
				ID:    "1",
				First: "allee",
				Last:  "crabb",
				Email: "crabballee@gmail.com",
			},
			{
				ID:    "2",
				First: "evan",
				Last:  "lurvey",
				Phone: "417-576-1238",
			},
		},
	}
	return db
}

func (db *DB) FindContacts() []Contact {
	db.m.RLock()
	defer db.m.RUnlock()
	return db.contacts
}

func (db *DB) GetContactByID(id string) (Contact, bool) {
	db.m.RLock()
	defer db.m.RUnlock()
	for _, c := range db.contacts {
		if c.ID == id {
			return c, true
		}
	}
	return Contact{}, false
}

func (db *DB) AddContact(in Contact) {
	db.m.Lock()
	defer db.m.Unlock()
	db.contacts = append(db.contacts, in)
}

func (db *DB) UpdateContactByID(in Contact) bool {
	db.m.Lock()
	defer db.m.Unlock()
	for i, c := range db.contacts {
		if c.ID == in.ID {
			db.contacts[i] = in
			return true
		}
	}
	return false
}
