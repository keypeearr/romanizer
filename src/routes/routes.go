package routes

import (
	"github.com/gofiber/fiber/v3"

	romanizerhandlers "github.com/keypeearr/romanizer/src/handlers/romanizerHandlers"
)

func Load(app *fiber.App) {
	views := app.Group("")
	views.Get("/", func(ctx fiber.Ctx) error {
		return ctx.Redirect().Status(fiber.StatusMovedPermanently).To("/romanizer")
	})
	views.Get("/romanizer", romanizerhandlers.DisplayRomanizer)

	api := app.Group("/api/v1/romanizer")
	api.Post("/alphaToRoman", romanizerhandlers.DisplayRomanValue)
	api.Post("/romanToAlpha", romanizerhandlers.DisplayAlphaValue)
}
