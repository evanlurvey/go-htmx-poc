package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// TODO: current thought is I need to make a form input field struct that can be passed in
// to auto build forms and handle error responses and what not in a standard way.
// should be helpful to keep things consistent and quick.

var appversion string
var csrfSecret string = "secret"

func NewCSRFToken(sessionID string) string {
	mac := hmac.New(sha256.New, []byte(csrfSecret))
	_, _ = mac.Write([]byte(sessionID))
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyCSRFToken(sessionID, token string) bool {
	og, err := hex.DecodeString(token)
	if err != nil {
		return false
	}
	// PERF: I know this is wasteful but whatever
	expected, _ := hex.DecodeString(NewCSRFToken(sessionID))
	return hmac.Equal(expected, og)
}

func init() {
	appversion = newID()
}


//go:embed views
var templatesFS embed.FS

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

func (TemplateEngine) dict(args ...any) (map[string]any, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]any, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = args[i+1]
	}
	return dict, nil
}

func (e *TemplateEngine) Render(c *fiber.Ctx, name string, binding interface{}, layout ...string) error {
	t := template.New("__root__").Funcs(template.FuncMap{
		"dict":       e.dict,
		"AppVersion": func() string { return appversion },
		"CSRFToken":  func() string { return NewCSRFToken(c.Cookies("session")) },
		"CSRFTokenInput": func() template.HTML {
			token := NewCSRFToken(c.Cookies("session"))
			return template.HTML(`<input type="hidden" name="csrf_token" value="` + token + `" /> `)
		},
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

type controller interface {
	Setup(fiber.Router)
}

func setupController(app fiber.Router, c controller) {
	c.Setup(app)
}

func main() {
	db := &DB{
		contacts: []Contact{
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
		},
	}

	engine := &TemplateEngine{
		fs:            templatesFS,
		defaultLayout: "layouts/main.html",
	}

	app := fiber.New(fiber.Config{
		Immutable:             true, // string buffers get reused otherwise and shit gets weird when using in memory things
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

	setupController(app, &ContactsRouter{
		templates: engine,
		db:        db,
	})

	slog.Info("starting server")
	if err := app.Listen(":8080"); err != nil {
		slog.Error("failed to start server", err)
	}
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

func newID() string {
	var b = make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
