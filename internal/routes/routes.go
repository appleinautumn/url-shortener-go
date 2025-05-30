package routes

import (
	"net/http"

	"url-shortener-go/internal/handlers"

	"github.com/go-chi/chi/v5"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/shorten", handlers.CreateShortURL)
	r.Get("/{shortID}", handlers.GetURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener"))
	})

	return r
}
