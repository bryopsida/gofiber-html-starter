package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// AddSwagger Registers the swagger ui page with the fiber app
// - app: *fiber.App fiber app
func AddSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
