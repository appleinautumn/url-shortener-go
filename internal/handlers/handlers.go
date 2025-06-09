package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"url-shortener-go/internal/services"

	"github.com/go-chi/chi/v5"
)

type URLMapping struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

type Handler struct {
	service services.URLService
}

func NewHandler(service services.URLService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var input URLMapping
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Invalid JSON in request body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	short, err := h.service.CreateShortURL(input.Long)
	if err != nil {
		slog.Error("Failed to create short URL", "error", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	slog.Info("URL created", "short", short, "long", input.Long)
	json.NewEncoder(w).Encode(URLMapping{Short: short, Long: input.Long})
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")

	longURL, err := h.service.GetLongURL(shortID)
	if err == sql.ErrNoRows {
		slog.Warn("URL not found", "short", shortID)
		http.NotFound(w, r)
		return
	} else if err != nil {
		slog.Error("Failed to retrieve URL", "error", err, "short", shortID)
		http.Error(w, "Failed to retrieve URL", http.StatusInternalServerError)
		return
	}

	slog.Info("URL retrieved", "short", shortID, "long", longURL)
	json.NewEncoder(w).Encode(URLMapping{Short: shortID, Long: longURL})
}
