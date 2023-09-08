package template

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"
	"htmx-poc/app"
	"maps"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type HTML = template.HTML

//go:embed views
var templatesFS embed.FS

var templateFuncs = map[string]func(ctx context.Context) any{}

func RegisterComponent(name string, c func(ctx context.Context) any) {
	// validate
	o := reflect.TypeOf(c(context.Background())).Kind()
	if o != reflect.Func {
		panic("output should be a function")
	}
	if _, exists := templateFuncs[name]; exists {
		panic("you cannot overwrite other functions")
	}
	templateFuncs[name] = c
}

type FS interface {
	ReadFile(string) ([]byte, error)
}

type TemplateEngine struct {
	fs            FS
	defaultLayout string
}

func NewTemplateEngine(fs FS, defaultLayout string) TemplateEngine {
	if fs == nil {
		fs = templatesFS
	}
	return TemplateEngine{
		fs:            fs,
		defaultLayout: defaultLayout,
	}
}
func (e TemplateEngine) RenderComponent(ctx context.Context, name string, data map[string]any) (template.HTML, error) {
	// main template
	f, err := e.openFile(name)
	if err != nil {
		return "", err
	}
	t, err := template.New(name).Parse(f)
	if err != nil {
		return "", err
	}
	// magically add user to everything
	global := FromContext(ctx)
	maps.Copy(global, data) // binding will overwrite global on conflict

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, name, global)
	return template.HTML(buf.String()), err
}

func (e TemplateEngine) Render(c *fiber.Ctx, name string, data map[string]any, layouts ...string) error {
	ctx := c.UserContext()
	// TODO: Extract funcs into something else
	t := template.New("__root__").Funcs(template.FuncMap{
		"dict":       e.dict,
		"AppVersion": app.AppVersion,
	})

	// add all of template funcs as options to call
	fm := make(template.FuncMap, len(templateFuncs))
	for k, v := range templateFuncs {
		fm[k] = v(ctx)
	}
	t.Funcs(fm)

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
	// magically add user to everything
	global := FromContext(ctx)
	maps.Copy(global, data) // binding will overwrite global on conflict

	c.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	// slim response on htmx shiz
	if c.Get("HX-Request") != "" {
		return t.ExecuteTemplate(c, "content", global)
	}
	return t.ExecuteTemplate(c, name, global)
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
