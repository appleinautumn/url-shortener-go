package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"url-shortener-go/internal/config"
	"url-shortener-go/internal/routes"
	"url-shortener-go/storage"

	"log/slog"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("Environment", "APP_ENV", cfg.AppEnv)

	if err := storage.InitDB(cfg.DBFile); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer storage.CloseDB()

	r := routes.Routes()

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
