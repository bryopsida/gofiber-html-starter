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
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"

	"github.com/bryopsida/gofiber-pug-starter/config"
	"github.com/bryopsida/gofiber-pug-starter/datastore"
	_ "github.com/bryopsida/gofiber-pug-starter/docs"
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/bryopsida/gofiber-pug-starter/pages"
	"github.com/bryopsida/gofiber-pug-starter/repositories/number"
	incrementroutes "github.com/bryopsida/gofiber-pug-starter/routes/increment"
	"github.com/bryopsida/gofiber-pug-starter/services/increment"
)

//go:embed public/*
var embedDirPubic embed.FS

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

func attachMiddleware(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(etag.New())
	app.Use(requestid.New())
	app.Use(requestid.New())
	app.Use(cors.New())
	app.Use(csrf.New())
	app.Use(compress.New())
	app.Use(cache.New())
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

func main() {
	slog.Info("Starting")
	config := config.NewViperConfig()
	slog.Info("Getting database")
	db, err := datastore.GetDatabase(config)
	if err != nil {
		slog.Error("failed to get database", "error", err)
		panic(err.Error())
	}
	defer db.Close()

	slog.Info("Getting number repository")
	repo := number.NewBadgerNumberRepository(db)

	slog.Info("Getting increment service")
	service := increment.NewIncrementService(repo, "counter")

	// Create a context with cancellation
	_, cancel := context.WithCancel(context.Background())

	// ensure this is always called on func exit
	defer cancel()

	appViews := buildViewEngine()
	appConfig := buildConfig(appViews)
	app := buildApp(appConfig)
	attachMiddleware(app)

	slog.Info("Registering routes")
	incrementroutes.RegisterRoutes(app, service)

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
