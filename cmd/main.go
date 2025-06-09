package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"url-shortener-go/internal/config"
	"url-shortener-go/internal/handlers"
	"url-shortener-go/internal/repository"
	"url-shortener-go/internal/routes"
	"url-shortener-go/internal/services"
	"url-shortener-go/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	var level slog.Level
	switch cfg.AppEnv {
	case "production":
		level = slog.LevelWarn
	case "development":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
	slog.Info("Environment", "APP_ENV", cfg.AppEnv)

	// Initialize database
	db, err := storage.InitDB(cfg.DBFile)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repository, service, and handler
	urlRepo := repository.NewURLRepository(db)
	urlService := services.NewURLService(urlRepo)
	handler := handlers.NewHandler(urlService)

	r := routes.Routes(handler)

	slog.Info("Listening on port", "port", cfg.AppPort)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe error", "error", err)
		}
	}()

	// Block until a signal is received
	<-stop

	slog.Info("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", "error", err)
	}

	slog.Info("Server exited properly")
}
