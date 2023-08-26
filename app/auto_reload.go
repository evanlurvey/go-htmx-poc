package app

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupAutoReloadWS(app *fiber.App, appversion string) {
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
		if err = c.WriteMessage(websocket.TextMessage, []byte(appversion)); err != nil {
			fmt.Println("open err:", err)
		}
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				fmt.Println("read:", err)
				break
			}
			fmt.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				fmt.Println("write:", err)
				break
			}
		}

	}))
}
