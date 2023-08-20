package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log/slog"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var appversion string

func init() {
	var b = make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	appversion = hex.EncodeToString(b)
}

//go:embed views
var templatesFS embed.FS

type Contact struct {
	ID    string
	First string
	Last  string
	Phone string
	Email string
}

var contacts = []Contact{
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
}

type TemplateEngine struct {
	fs            embed.FS
	defaultLayout string
}

func (e *TemplateEngine) openFile(name string) string {
	ts, err := e.fs.ReadFile("views/" + name)
	if err != nil {
		panic(err)
	}
	return string(ts)
}

func (e *TemplateEngine) Render(c *fiber.Ctx, name string, binding interface{}, layout ...string) error {
	t := template.New("__root__").Funcs(template.FuncMap{
		"AppVersion": func() string { return appversion },
	})
	var err error
	// default layout load
	if e.defaultLayout != "" && len(layout) == 0 {
		t, err = t.New(e.defaultLayout).Parse(e.openFile(e.defaultLayout))
		if err != nil {
			return err
		}
	}
	// optional layouts load
	for _, l := range layout {
		t, err = t.New(l).Parse(e.openFile(l))
		if err != nil {
			return err
		}
	}
	// main template
	t, err = t.New(name).Parse(e.openFile(name))
	if err != nil {
		return err
	}
	c.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	return t.ExecuteTemplate(c, name, binding)

}

func setupAutoReloadWS(app *fiber.App) {
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

func main() {
	engine := &TemplateEngine{
		fs:            templatesFS,
		defaultLayout: "layouts/main.html",
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	// FIX: remove in prod
	setupAutoReloadWS(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return engine.Render(c, "index.html", map[string]any{
			"ctx":  c.UserContext(),
			"name": "evan<p>lol</p>",
		})
	})

	app.Route("/contacts", func(r fiber.Router) {
		r.Get("/", func(c *fiber.Ctx) error {
			return engine.Render(c, "contacts.html", map[string]any{
				"contacts": contacts,
			})
		})
		r.Get("/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			slog.Info("contact detail", slog.String("contact_id", id))
			contact, found := get(contacts, func(c Contact) bool { return c.ID == id })
			if !found {
				return c.SendStatus(404)
			}
			return engine.Render(c, "contacts-detail.html", map[string]any{
				"contact": contact,
			})
		})
		r.Get("/:id/edit", func(c *fiber.Ctx) error {
			id := c.Params("id")
			slog.Info("contact detail", slog.String("contact_id", id))
			contact, found := get(contacts, func(c Contact) bool { return c.ID == id })
			if !found {
				return c.SendStatus(404)
			}
			return engine.Render(c, "contacts-form.html", map[string]any{
				"contact": contact,
			})
		})
		r.Get("/new", func(c *fiber.Ctx) error {
			c.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
			return c.SendStatus(200)

		})
	})

	slog.Info("starting server")
	if err := app.Listen(":8080"); err != nil {
		slog.Error("failed to start server", err)
	}
}

func get[T any](arr []T, f func(v T) bool) (out T, found bool) {
	for _, v := range arr {
		if f(v) {
			slog.Info("found ya")
			return v, true
		}
	}
	return
}
