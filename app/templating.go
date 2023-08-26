package app

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"
	"htmx-poc/app/csrf"

	"github.com/gofiber/fiber/v2"
)

//go:embed views
var templatesFS embed.FS

type TemplateEngine struct {
	fs            embed.FS
	csrf          csrf.Service
	defaultLayout string
	appversion    string
}

func NewTemplateEngine(csrfService csrf.Service, appversion string, defaultLayout string) TemplateEngine {
	return TemplateEngine{
		fs:            templatesFS,
		csrf:          csrfService,
		defaultLayout: defaultLayout,
		appversion:    appversion,
	}
}

func (e TemplateEngine) Render(c *fiber.Ctx, name string, binding map[string]any, layouts ...string) error {
	ctx := c.UserContext()
	// TODO: Extract funcs into something else
	t := template.New("__root__").Funcs(template.FuncMap{
		"dict":       e.dict,
		"AppVersion": e.appVersion,
		"form":       e.form(ctx),
	})

	var err error
	e.loadLayout(t, layouts...)

	// main template
	t, err = t.New(name).Parse(e.openFile(name))
	if err != nil {
		return err
	}
	c.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	return t.ExecuteTemplate(c, name, binding)
}

func (e TemplateEngine) openFile(name string) string {
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

func (e TemplateEngine) appVersion() string { return e.appversion }

func (e TemplateEngine) form(ctx context.Context) func(f Form) (template.HTML, error) {
	return func(f Form) (template.HTML, error) {
		t, err := template.New("form").Parse(e.openFile("components/form.html"))
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		err = t.Execute(&buf, map[string]any{
			"form": f,
		})
		return template.HTML(buf.String()), err
	}
}

func (e TemplateEngine) loadLayout(t *template.Template, layouts ...string) error {
	var err error
	// default layout load
	if e.defaultLayout != "" && len(layouts) == 0 {
		t, err = t.New(e.defaultLayout).Parse(e.openFile(e.defaultLayout))
		if err != nil {
			return err
		}
		return nil
	}
	// optional layouts load
	for _, l := range layouts {
		t, err = t.New(l).Parse(e.openFile(l))
		if err != nil {
			return err
		}
	}
	return nil
}
