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

func (a *authRoutes) LoginHandler(c *fiber.Ctx) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")
	slog.Info("Login attempt for user", "user", user)
	dbUser, err := a.userService.GetUserByUsername(user)
	if err != nil {
		slog.Info("Failed login attempt for user", "user", user)
		c.Redirect("/login?loginError=true")
		return nil
	}
	validPass, err := a.passwordService.Verify(pass, dbUser.PasswordHash)
	if err != nil || !validPass {
		slog.Info("Invalid login credentials provided for user", "user", user)
		c.Redirect("/login?loginError=true")
		return nil
	}
	// Create the Claims
	token, err := a.jwtService.Generate(dbUser)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		c.Redirect("/login?loginError=true")
		return nil
	}
	// stick the token in the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "app_user",
		Value:    token,
		SameSite: "Strict",
		HTTPOnly: true,
	})
	c.Redirect("/")
	return nil
}

func (a *authRoutes) LogoutHandler(c *fiber.Ctx) error {
	return nil
}

func RegisterPublicRoutes(router fiber.Router, passwordService interfaces.IPasswordService, userService interfaces.IUsersService, jwtService interfaces.IJWTService) {
	slog.Info("Adding public auth routes", "router", router)
	authRoutes := authRoutes{passwordService: passwordService, userService: userService, jwtService: jwtService}

	router.Post("/login", authRoutes.LoginHandler)

}

func RegisterPrivateRoutes(router fiber.Router, passwordService interfaces.IPasswordService, userService interfaces.IUsersService, jwtService interfaces.IJWTService) {
	slog.Info("Adding private auth routes", "router", router)
	authRoutes := authRoutes{passwordService: passwordService, userService: userService, jwtService: jwtService}

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
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			slog.Error("JWT Error", "error", err)
			return c.Redirect("/login")
		},
		TokenLookup: "cookie:app_user",
	}))
}
