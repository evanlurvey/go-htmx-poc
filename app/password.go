package app

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PHC struct {
	algorithm string
	version   int
	memory    uint32
	time      uint32
	threads   uint8
	salt      []byte
	hash      []byte
}

func (phc PHC) String() string {
	const format = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, phc.memory, phc.time, phc.threads, hex.EncodeToString(phc.salt), hex.EncodeToString(phc.hash))
}

func (phc PHC) VerifyPassword(password string) bool {
	var keyLen = uint32(len(phc.hash))

	comparisonHash := argon2.IDKey([]byte(password), phc.salt, phc.time, phc.memory, phc.threads, keyLen)

	return subtle.ConstantTimeCompare(phc.hash, comparisonHash) == 1
}

func phcFromDB(s string) PHC {
	var phc PHC
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

func NewPHC(password string) PHC {
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
	return PHC{
		algorithm: "argon2id",
		version:   argon2.Version,
		memory:    memory,
		time:      time,
		threads:   threads,
		salt:      salt,
		hash:      hash,
	}
}
