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
	fs            interface{ ReadFile(string) ([]byte, error) }
	csrf          csrf.Service
	defaultLayout string
}

func NewTemplateEngine(csrfService csrf.Service, defaultLayout string) TemplateEngine {
	return TemplateEngine{
		fs:            templatesFS,
		csrf:          csrfService,
		defaultLayout: defaultLayout,
	}
}

func (e TemplateEngine) Render(c *fiber.Ctx, name string, binding map[string]any, layouts ...string) error {
	ctx := c.UserContext()
	// TODO: Extract funcs into something else
	t := template.New("__root__").Funcs(template.FuncMap{
		"dict":       e.dict,
		"AppVersion": AppVersion,
		"form":       e.form(ctx),
	})

	var err error
	e.loadLayout(t, layouts...)

	// main template
	f, err := e.openFile(name)
	if err != nil {
		return err
	}
	t, err = t.New(name).Parse(f)
	if err != nil {
		return err
	}
	c.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	return t.ExecuteTemplate(c, name, binding)
}

func (e TemplateEngine) openFile(name string) (string, error) {
	ts, err := e.fs.ReadFile("views/" + name)
	if err != nil {
		return "", err
	}
	return string(ts), nil
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

func (e TemplateEngine) form(ctx context.Context) func(f Form) (template.HTML, error) {
	return func(f Form) (template.HTML, error) {
		f = f.AddCSRFToken(ctx, e.csrf)
		if f.Template == "" {
			f.Template = "components/form.html"
		}
		tf, err := e.openFile(f.Template)
		if err != nil {
			return "", err
		}

		t, err := template.New("form").Parse(tf)
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
	// default layout load
	if e.defaultLayout != "" && len(layouts) == 0 {
		f, err := e.openFile(e.defaultLayout)
		if err != nil {
			return err
		}
		t, err = t.New(e.defaultLayout).Parse(f)
		if err != nil {
			return err
		}
		return nil
	}
	// optional layouts load
	for _, l := range layouts {
		f, err := e.openFile(l)
		if err != nil {
			return err
		}
		t, err = t.New(l).Parse(f)
		if err != nil {
			return err
		}
	}
	return nil
}
