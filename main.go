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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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
	password_service "github.com/bryopsida/gofiber-pug-starter/services/password"
	settings_service "github.com/bryopsida/gofiber-pug-starter/services/settings"

	incrementroutes "github.com/bryopsida/gofiber-pug-starter/routes/increment"
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
}

func buildConfig(view fiber.Views) fiber.Config {
	return fiber.Config{
		Views:       view,
		ViewsLayout: "layouts/main",
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

	app.Use(helmet.New())
	app.Use(etag.New())
	app.Use(requestid.New())
	app.Use(requestid.New())
	app.Use(cors.New())
	app.Use(csrf.New())
	app.Use(compress.New())
	app.Use(cache.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptionKey,
	}))
	app.Use(logger.New())
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
	return services
}

func main() {
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

	slog.Info("Registering routes")
	incrementroutes.RegisterRoutes(app, services.IncrementService)

	slog.Info("Registering global pages")
	pages.RegisterGlobalPages(app)

	slog.Info("Adding Swagger")
	pages.AddSwagger(app)

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
