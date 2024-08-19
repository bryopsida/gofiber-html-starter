package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug/v2"

	"github.com/bryopsida/gofiber-pug-starter/config"
	"github.com/bryopsida/gofiber-pug-starter/datastore"
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/bryopsida/gofiber-pug-starter/pages"
	"github.com/bryopsida/gofiber-pug-starter/repositories/number"
	incrementroutes "github.com/bryopsida/gofiber-pug-starter/routes/increment"
	"github.com/bryopsida/gofiber-pug-starter/services/increment"
)

func buildConfig(view *pug.Engine) fiber.Config {
	return fiber.Config{
		Views: view,
	}
}

func buildViewEngine() *pug.Engine {
	engine := pug.New("./views", ".pug")
	engine.Load()
	return engine
}

func buildApp(config fiber.Config) *fiber.App {
	return fiber.New(config)
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

	slog.Info("Registering routes")
	incrementroutes.RegisterRoutes(app, service)

	slog.Info("Registering global pages")
	pages.RegisterGlobalPages(app)

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
