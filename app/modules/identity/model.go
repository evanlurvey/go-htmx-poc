package identity

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"htmx-poc/utils"
	"io"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

type user struct {
	User
	password phc
}

type loginOutcome uint8

const (
	loginOutcome_nil loginOutcome = iota
	loginOutcome_success
	loginOutcome_invalidEmail
	loginOutcome_invalidPassword
)

type loginAttempts []loginAttempt

func (la loginAttempts) Unsuccessful() int {
	return len(
		utils.FilterFunc(
			la,
			func(la loginAttempt) bool { return la.outcome != loginOutcome_success },
		),
	)
}

type loginAttempt struct {
	id      string
	at      time.Time
	user_id string // optional field
	email   string
	outcome loginOutcome
}

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

func (s Session) ID() string {
	return s.id
}

func (s Session) User() (User, bool) {
	if s.state == SessionState_Authenticated && s.Valid() {
		return s.user, true
	}
	return User{}, false
}

func (s Session) Valid() bool {
	return time.Now().Before(s.expires)
}

type phc struct {
	algorithm string
	version   int
	memory    uint32
	time      uint32
	threads   uint8
	salt      []byte
	hash      []byte
}

func (p phc) String() string {
	const format = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, p.memory, p.time, p.threads, hex.EncodeToString(p.salt), hex.EncodeToString(p.hash))
}

func (p phc) VerifyPassword(password string) bool {
	var keyLen = uint32(len(p.hash))

	comparisonHash := argon2.IDKey([]byte(password), p.salt, p.time, p.memory, p.threads, keyLen)

	return subtle.ConstantTimeCompare(p.hash, comparisonHash) == 1
}

func phcFromDB(s string) phc {
	var phc phc
	var hashStr string
	var saltStr string
	var err error

	s = strings.TrimSpace(strings.ReplaceAll(s, "$", " "))
	_, err = fmt.Sscanf(s, "%s v=%d m=%d,t=%d,p=%d %s %s", &phc.algorithm, &phc.version, &phc.memory, &phc.time, &phc.threads, &saltStr, &hashStr)
	if err != nil {
		panic(err)
	}

	phc.hash, err = hex.DecodeString(hashStr)
	if err != nil {
		panic("failed to decode hash, this should never happen.")
	}

	phc.salt, err = hex.DecodeString(saltStr)
	if err != nil {
		panic("failed to decode salt, this should never happen.")
	}

	return phc
}

func newPHC(password string) phc {
	// time is iterations to perform
	const time = 3
	const saltLen = 16
	// memory in KiB
	const memory = 64 * 1024
	const threads = 4
	const keyLen = 64

	salt := make([]byte, saltLen)
	if l, _ := io.ReadFull(rand.Reader, salt); l != saltLen {
		panic("couldn't fill salt buffer from random reader, this should never happen.")
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	return phc{
		algorithm: "argon2id",
		version:   argon2.Version,
		memory:    memory,
		time:      time,
		threads:   threads,
		salt:      salt,
		hash:      hash,
	}
}
