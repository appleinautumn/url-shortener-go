package routes

import (
	"net/http"

	"url-shortener-go/internal/handlers"
	"url-shortener-go/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middlewares.JSONContentType)
	r.Use(middlewares.Logging)

	r.Post("/shorten", handler.CreateURL)
	r.Get("/{shortID}", handler.GetURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
