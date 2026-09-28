package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

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
	
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalUrl := strings.TrimSpace(string(bytes))
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
	w.Write([]byte(h.baseUrl + "/" + shortUrl))
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
	w.Write([]byte(originalUrl))
}
