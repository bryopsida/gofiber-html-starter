package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func AddSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
