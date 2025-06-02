package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"url-shortener-go/internal/routes"
	"url-shortener-go/storage"

	"github.com/joho/godotenv"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		log.Fatalf("DB_FILE must be set in the .env file")
	}

	if err := storage.InitDB(dbFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.CloseDB()

	r := routes.Routes()

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Listening on :%s...", appPort)

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
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Block until a signal is received
	<-stop

	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}
