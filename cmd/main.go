package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"url-shortener-go/storage"

	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Post("/shorten", createShortURL)
	r.Get("/{shortID}", getURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener"))
	})

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Listening on :%s...", appPort)
	http.ListenAndServe(":"+appPort, r)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var input URLMapping
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate random short ID
	short := fmt.Sprintf("%x", rand.IntN(100000))

	if err := storage.StoreURL(short, input.Long); err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(URLMapping{Short: short, Long: input.Long})
}

func getURL(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")

	var longURL string
	longURL, err := storage.GetLongURL(shortID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(URLMapping{Short: shortID, Long: longURL})
}
