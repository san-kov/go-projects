package handler

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"url-shortener/models"
	"url-shortener/storage"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Code string `json:"code"`
}

type StatsResponse struct {
	URL    string `json:"url"`
	Code   string `json:"code"`
	Clicks int    `json:"clicks"`
}

func randomCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

type Handler struct {
	storage storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
	return &Handler{storage: s}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	code := randomCode(6)

	if err := h.storage.Save(models.NewURLEntry(code, req.URL)); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{Code: code})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	entry, err := h.storage.Get(code)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "code not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	entry.Click()
	http.Redirect(w, r, entry.OriginalURL, http.StatusMovedPermanently)

}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/stats/"):]

	entry, err := h.storage.Get(code)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "code not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatsResponse{URL: entry.OriginalURL, Code: entry.Code, Clicks: entry.Clicks})
}
