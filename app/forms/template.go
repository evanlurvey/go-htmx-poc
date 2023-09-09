package forms

import (
	"context"
	"embed"
	"htmx-poc/app/template"
)

//go:embed views
var templatesFS embed.FS
var templateEngine = template.NewTemplateEngine(templatesFS, "")

func form(ctx context.Context) any {
	return func(f Form) (template.HTML, error) {
		f = f.AddCSRFToken(ctx)
		if f.Template == "" {
			f.Template = "form.html"
		}
		return templateEngine.RenderComponentHTML(ctx, f.Template, template.M{
			"form": f,
		})
	}
}

func init() {
	template.RegisterComponent("form", form)
}
