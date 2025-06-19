package utils

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func Render(ctx fiber.Ctx, component templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}
