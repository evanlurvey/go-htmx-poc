package web

import (
	"embed"
	"io"
	"mime"
	"net/http"
	"path"

	"github.com/gofiber/fiber/v2"
)

//go:embed static
var staticFS embed.FS

func setupStatic(web *fiber.App) {
	web.Get("/static/+", func(c *fiber.Ctx) error {
		p := "static/" + c.Params("+")
		f, err := staticFS.Open(p)
		if err != nil {
			return c.SendStatus(404)
		}
		s, _ := f.Stat()
		ext := path.Ext(s.Name())
		contentType := mime.TypeByExtension(ext)
		if len(contentType) == 0 {
			var data [512]byte
			_, err := io.ReadFull(f, data[:])
			if err != nil {
				return err
			}
			contentType = http.DetectContentType(data[:])
		}
		c.Set("Content-Type", contentType)
		// TODO: Caching headers

		// reading whole file instead of sending as stream so
		// compression middleware will do our job for us
		o, _ := io.ReadAll(f)
		return c.Send(o)
	})
}
