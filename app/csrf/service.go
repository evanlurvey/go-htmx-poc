package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/gofiber/fiber/v2"
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

func (c Service) NewToken(sessionID string) string {
	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write([]byte(sessionID))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c Service) VerifyToken(sessionID, token string) error {
	og, err := hex.DecodeString(token)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write([]byte(sessionID))
	if !hmac.Equal(mac.Sum(nil), og) {
		return InvalidCSRFError
	}
	return nil
}

func (Service) ErrorHandler(c *fiber.Ctx) error {
	err := c.Next()
	if errors.Is(err, InvalidCSRFError) {
		return c.Status(401).SendString("invalid csrf")
	}
	return err
}
