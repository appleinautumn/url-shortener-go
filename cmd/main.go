package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"url-shortener-go/internal/routes"
	"url-shortener-go/storage"

	"log/slog"

	"github.com/joho/godotenv"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	// Get environment
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	slog.Info("Environment", "APP_ENV", appEnv)

	// Get database location
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		slog.Error("DB_FILE must be set in the .env file")
		os.Exit(1)
	}

	if err := storage.InitDB(dbFile); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer storage.CloseDB()

	r := routes.Routes()

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	slog.Info("Listening on port", "port", appPort)

	srv := &http.Server{
		Addr:    ":" + appPort,
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
