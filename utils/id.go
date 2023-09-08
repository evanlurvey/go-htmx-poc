package utils

import (
	"crypto/rand"
	"encoding/base64"
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
	return base64.URLEncoding.EncodeToString(b)
}
