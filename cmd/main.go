package main

import (
	"log"
	"net/http"
	"os"

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
	http.ListenAndServe(":"+appPort, r)
}
