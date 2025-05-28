package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"

	"url-shortener-go/storage"

	"github.com/go-chi/chi/v5"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/shorten", createShortURL)
	r.Get("/{shortID}", getURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener"))
	})

	return r
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
