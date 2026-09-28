package handler

import (
	"errors"
	"net/http"
	"strings"
	"encoding/json"

	"short-urls/internal/service"
	"short-urls/internal/storage"
)

type Handler struct {
	service *service.Service
	baseUrl string
}

func NewHandler(service *service.Service, baseUrl string) *Handler {
	baseUrl = strings.TrimRight(baseUrl, "/")
	return &Handler{service: service, baseUrl: baseUrl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /generate", h.HandlerGenerateShortUrl)
	mux.HandleFunc("GET /", h.HandlerGetOriginalUrl)
}

func (h *Handler) HandlerGenerateShortUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RequestOriginalUrl
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalUrl := strings.TrimSpace(req.OriginalUrl)
	if originalUrl == "" {
		http.Error(w, "Empty url", http.StatusBadRequest)
		return
	}

	shortUrl, err := h.service.GenerateShortUrl(ctx, originalUrl)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	responseJson := ResponseShortUrl{ShortUrl: h.baseUrl + "/" + shortUrl}
	if err = json.NewEncoder(w).Encode(responseJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandlerGetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	

	shortUrl := strings.TrimPrefix(r.URL.Path, "/")
	if shortUrl == "" {
		http.Error(w, "ShortUrl should not empty", http.StatusBadRequest)
		return
	}

	originalUrl, err := h.service.GetOriginalUrl(ctx, shortUrl)
	if err != nil {
		if errors.Is(err, storage.ErrUrlNotFound) {
			http.Error(w, "ShortUrl not found", http.StatusNotFound)
			return	
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	responseJson := ResponseOriginalUrl{OriginalUrl: originalUrl}
	if err = json.NewEncoder(w).Encode(responseJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type ResponseShortUrl struct {
	ShortUrl string `json:"short_url"`
}

type ResponseOriginalUrl struct {
	OriginalUrl string `json:"original_url"`
}

type RequestOriginalUrl struct {
	OriginalUrl string `json:"original_url"`
}
