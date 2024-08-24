package pages

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

// RegisterGlobalPages registers global pages
// - app: *fiber.App fiber app
func RegisterGlobalPages(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		slog.Info("Rendering index")
		return c.Render("index", fiber.Map{
			"cardRows": []fiber.Map{},
		})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		slog.Info("Rendering login")
		return c.Render("login", fiber.Map{})
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		slog.Info("Rendering about")
		return c.Render("about", fiber.Map{})
	})

	app.Get("/404", func(c *fiber.Ctx) error {
		slog.Info("Rendering 404")
		return c.Render("404", fiber.Map{})
	})

	app.Get("/500", func(c *fiber.Ctx) error {
		slog.Info("Rendering 500")
		return c.Render("500", fiber.Map{})
	})

	app.Get("/403", func(c *fiber.Ctx) error {
		slog.Info("Rendering 403")
		return c.Render("403", fiber.Map{})
	})

	app.Get("/401", func(c *fiber.Ctx) error {
		slog.Info("Rendering 401")
		return c.Render("401", fiber.Map{})
	})
}
