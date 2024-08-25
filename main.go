package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html/v2"
	"github.com/pressly/goose/v3"
	slogfiber "github.com/samber/slog-fiber"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bryopsida/gofiber-pug-starter/auth"
	"github.com/bryopsida/gofiber-pug-starter/config"
	"github.com/bryopsida/gofiber-pug-starter/database"
	"github.com/bryopsida/gofiber-pug-starter/database/migrations"
	_ "github.com/bryopsida/gofiber-pug-starter/docs"
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/bryopsida/gofiber-pug-starter/pages"
	number_repsitory "github.com/bryopsida/gofiber-pug-starter/repositories/number"
	settings_repository "github.com/bryopsida/gofiber-pug-starter/repositories/settings"
	users_repository "github.com/bryopsida/gofiber-pug-starter/repositories/users"
	increment_service "github.com/bryopsida/gofiber-pug-starter/services/increment"
	jwt_service "github.com/bryopsida/gofiber-pug-starter/services/jwt"
	password_service "github.com/bryopsida/gofiber-pug-starter/services/password"
	settings_service "github.com/bryopsida/gofiber-pug-starter/services/settings"
	users_service "github.com/bryopsida/gofiber-pug-starter/services/users"
)

//go:embed database/migrations/sql/*
var sqlMigrations embed.FS

//go:embed public/*
var embedDirPubic embed.FS

type repositories struct {
	NumberRepository   interfaces.INumberRepository
	SettingsRepository interfaces.ISettingsRepository
	UsersRepository    interfaces.IUserRepository
}

type services struct {
	IncrementService interfaces.IIncrementService
	SettingsService  interfaces.ISettingsService
	PasswordService  interfaces.IPasswordService
	UsersService     interfaces.IUsersService
	JWTService       interfaces.IJWTService
}

func buildConfig(view fiber.Views) fiber.Config {
	return fiber.Config{
		Views:                 view,
		ViewsLayout:           "layouts/main",
		PassLocalsToViews:     true,
		DisableStartupMessage: true,
	}
}

func buildViewEngine() *html.Engine {
	engine := html.New("./views", ".html")
	return engine
}

func buildApp(config fiber.Config) *fiber.App {
	return fiber.New(config)
}

func attachMiddleware(app *fiber.App, services *services) {
	// get the cookie encryption key
	encryptionKey, err := services.SettingsService.GetString("cookie_encryption_key")
	if err != nil {
		slog.Error("Error getting cookie encryption key", "error", err)
		panic("failed to get cookie encryption key")
	}
	app.Use(slogfiber.New(slog.Default()))
	app.Use(helmet.New())
	app.Use(etag.New())
	app.Use(requestid.New())
	app.Use(cors.New())
	app.Use(csrf.New(csrf.Config{
		KeyLookup:         "cookie:csrf_",
		CookieName:        "csrf_",
		CookieSameSite:    "Strict",
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Extractor:         csrf.CsrfFromCookie("csrf_"),
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
	}))
	app.Use(compress.New())
	app.Use(cache.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptionKey,
	}))
	app.Use(healthcheck.New())

	app.Use("/public", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirPubic),
		PathPrefix: "public",
		Browse:     false,
	}))

}

func runServer(app *fiber.App, address string) {
	err := app.Listen(address)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}

func startServer(app *fiber.App, config interfaces.IConfig) {
	address := config.GetServerAddress()
	port := config.GetServerPort()
	slog.Info("Starting server", "address", address, "port", port)
	serverListenAddress := fmt.Sprintf("%s:%d", address, port)
	go runServer(app, serverListenAddress)
}

func initializeDatabase(cfg interfaces.IConfig) *gorm.DB {

	var err error
	database.DBConn, err = gorm.Open(sqlite.Open(cfg.GetDatabasePath()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	slog.Info("Connection Opened to Database")
	sqlDb, err := database.DBConn.DB()
	if err != nil {
		panic("failed to get database connection")
	}
	goose.SetBaseFS(sqlMigrations)
	goose.SetDialect("sqlite3")
	migrations.InitializeV001Migration(*database.DBConn)
	migrations.InitializeV002Migration(*database.DBConn)
	migrations.InitializeV003Migration(*database.DBConn)
	migrations.InitializeV004Migration(*database.DBConn)
	err = goose.Up(sqlDb, "database/migrations/sql")

	if err != nil {
		slog.Error("Error migrating database", "error", err)
		panic("failed to migrate database")
	}
	slog.Info("Database Migrated")
	return database.DBConn
}

func initializeRepositories(db *gorm.DB) *repositories {
	// Initialize repositories
	repositories := &repositories{}
	repositories.NumberRepository = number_repsitory.NewNumberRepository(db)
	repositories.SettingsRepository = settings_repository.NewSettingsRepository(db)
	repositories.UsersRepository = users_repository.NewUserRepository(db)
	return repositories
}

func initializeServices(repos *repositories) *services {
	// Initialize services
	services := &services{}
	services.IncrementService = increment_service.NewIncrementService(repos.NumberRepository, "counter")
	services.PasswordService = password_service.NewPasswordService()
	services.SettingsService = settings_service.NewSettingsService(repos.SettingsRepository)
	services.JWTService = jwt_service.NewJWTService(services.SettingsService)
	services.UsersService = users_service.NewUsersService(repos.UsersRepository)
	return services
}

func addPublicRoutes(app *fiber.App, services *services) {
	apiV1Router := app.Group("/api/v1")
	auth.RegisterPublicRoutes(apiV1Router.Group("/auth"), services.PasswordService, services.UsersService, services.JWTService)
}
func addPublicPages(app *fiber.App) {
	pages.RegisterGlobalPages(app)
	pages.AddSwagger(app)
}

func addPrivateRoutes(app *fiber.App, services *services) {
	apiV1Router := app.Group("/api/v1")
	auth.RegisterPrivateRoutes(apiV1Router.Group("/auth"), services.PasswordService, services.UsersService, services.JWTService)
}
func addPrivatePages(app *fiber.App, services *services) {
	pages.RegisterPrivateGlobalPages(app, services.JWTService)
}

func addAuthMiddleware(app *fiber.App, services *services) {
	auth.AddJWTAuth(app, services.SettingsService)
}
func main() {
	defaultLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(defaultLogger)
	slog.Info("Starting")
	config := config.NewViperConfig()
	slog.Info("Getting database")
	db := initializeDatabase(config)

	repos := initializeRepositories(db)
	services := initializeServices(repos)

	// Create a context with cancellation
	_, cancel := context.WithCancel(context.Background())
	// ensure this is always called on func exit
	defer cancel()

	appViews := buildViewEngine()
	appConfig := buildConfig(appViews)
	app := buildApp(appConfig)
	attachMiddleware(app, services)
	addPublicRoutes(app, services)
	addPublicPages(app)
	addAuthMiddleware(app, services)
	addPrivateRoutes(app, services)
	addPrivatePages(app, services)

	startServer(app, config)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Wait for a signal
	sig := <-sigChan
	slog.Info("Received signal", "signal", sig)
	// Cancel the context
	cancel()
	slog.Info("Server stopped")
}
