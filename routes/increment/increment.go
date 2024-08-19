package incrementroutes

import (
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers increment routes
// - app: *fiber.App fiber app
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
