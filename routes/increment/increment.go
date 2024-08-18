package increment_routes

import (
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, service interfaces.IIncrementService) {

	app.Post("/increment", func(c *fiber.Ctx) error {
		id := c.Query("id")
		number, err := service.Increment(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"number": number})
	})
}
