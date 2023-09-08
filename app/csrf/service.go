package csrf

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var InvalidCSRFError = errors.New("invalid csrf")

type Service struct {
	secret []byte
}

func New(secret []byte) Service {
	return Service{
		secret: secret,
	}
}

func (c Service) NewToken(sid string) string {
	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write([]byte(sid))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func VerifyToken(ctx context.Context, token string) error {
	expected := FromContext(ctx)
	if token != expected {
		return InvalidCSRFError
	}
	return nil
}
