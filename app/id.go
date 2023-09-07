package app

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// size in bytes
func NewID(size ...uint8) string {
	var s uint8 = 16
	if len(size) == 1 {
		s = size[0]
	}
	var b = make([]byte, s)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
