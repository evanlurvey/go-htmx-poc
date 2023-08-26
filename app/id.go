package app

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func NewID() string {
	var b = make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
