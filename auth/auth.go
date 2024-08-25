package auth

import (
	"log/slog"

	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	passwordService interfaces.IPasswordService
	jwtService      interfaces.IJWTService
	userService     interfaces.IUsersService
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Login
// @Description Logs in a user and returns a JWT token
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 500 {object} map[string]string "error"
// @Router /api/v1/auth/login [post]
func (a *authRoutes) LoginHandler(c *fiber.Ctx) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")
	slog.Info("Login attempt for user", "user", user)
	dbUser, err := a.userService.GetUserByUsername(user)
	if err != nil {
		slog.Info("Failed login attempt for user", "user", user)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	validPass, err := a.passwordService.Verify(pass, dbUser.PasswordHash)
	if err != nil || !validPass {
		slog.Info("Invalid login credentials provided for user", "user", user)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// Create the Claims
	token, err := a.jwtService.Generate(dbUser)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// stick the token in the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "app_user",
		Value:    token,
		SameSite: "Strict",
		HTTPOnly: true,
	})
	// return it as well so it can be supplied in header
	return c.JSON(fiber.Map{"token": token})
}

// @Summary Logout
// @Description Close session
// @Tags auth
// @Success 200 {string} string "ok"
// @Router /api/v1/auth/logout [post]
func (a *authRoutes) LogoutHandler(c *fiber.Ctx) error {
	return nil
}

func RegisterRoutes(router fiber.Router, passwordService interfaces.IPasswordService, userService interfaces.IUsersService, jwtService interfaces.IJWTService) {
	slog.Info("Adding auth routes", "router", router)
	authRoutes := authRoutes{passwordService: passwordService, userService: userService, jwtService: jwtService}

	router.Post("/login", authRoutes.LoginHandler)
	router.Post("/logout", authRoutes.LogoutHandler)

}

func AddJWTAuth(app *fiber.App, settingsService interfaces.ISettingsService) {
	signingKey, err := settingsService.GetString("jwt_signing_key")
	if err != nil {
		slog.Error("Error getting JWT signing key", "error", err)
		panic(err.Error())
	}
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey)},
	}))
}
