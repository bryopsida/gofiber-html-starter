package pages

import (
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/gofiber/fiber/v2"
)

func validateAddUserForm(c *fiber.Ctx) (bool, error) {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirmPassword")

	if username == "" || email == "" || password == "" || confirmPassword == "" {
		return false, c.SendStatus(fiber.StatusBadRequest)
	}

	if password != confirmPassword {
		return false, c.SendStatus(fiber.StatusBadRequest)
	}

	return true, nil
}

func RegisterPrivateUserPages(app *fiber.App, userService interfaces.IUsersService, passwordService interfaces.IPasswordService) {
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.Render("users", fiber.Map{})
	})
	app.Get("/add-user", func(c *fiber.Ctx) error {
		return c.Render("add-user", fiber.Map{})
	})
	app.Post("/add-user", func(c *fiber.Ctx) error {
		valid, err := validateAddUserForm(c)
		if !valid {
			return err
		}

		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")
		role := c.FormValue("role")

		passwordHash, err := passwordService.Hash(password)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		err = userService.CreateUser(&interfaces.User{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
			Role:         role,
		})
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Redirect("/users")
	})
}
