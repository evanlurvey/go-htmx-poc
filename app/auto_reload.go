package app

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupAutoReloadWS(app *fiber.App) {
	app.Use("/dev/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("dev/ws/reload", websocket.New(func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		if err = c.WriteMessage(websocket.TextMessage, []byte(AppVersion())); err != nil {
			return
		}

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			}

			if err = c.WriteMessage(mt, msg); err != nil {
				break
			}
		}

	}))
}
