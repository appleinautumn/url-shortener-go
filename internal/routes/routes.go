package routes

import (
	"net/http"

	"url-shortener-go/internal/handlers"
	"url-shortener-go/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func Routes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middlewares.JSONContentType)

	r.Post("/shorten", handlers.CreateURL)
	r.Get("/{shortID}", handlers.GetURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
