package pages

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

// RegisterGlobalPages registers global pages
// - app: *fiber.App fiber app
func RegisterGlobalPages(app *fiber.App) {

	app.Get("/test", func(c *fiber.Ctx) error {
		slog.Info("Rendering test")
		err := c.Render("test", fiber.Map{})
		if err != nil {
			slog.Error("Error rendering test", "error", err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return nil
	})
}
