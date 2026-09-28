package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"short-urls/internal/config"
	"short-urls/internal/service"
	"short-urls/internal/storage"
)

func TestGenerateShortUrlAndGetOriginalUrl(t *testing.T) {
	cfg := config.Load()

	store := storage.NewMemoryStorage()
	svc := service.NewService(store)
	handler := NewHandler(svc, cfg.BaseURL)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	
	server := httptest.NewServer(mux)
	defer server.Close()

	originalUrl := "https://github.com"

	// post метод
	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(originalUrl),
	)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed generation to shortUrl: %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Waited status 201, but gor %d", response.StatusCode)
	}

	bytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	shortUrl := string(bytes)
	shortId := shortUrl[len(shortUrl)-service.LenghtShortUrl:]

	// get метод
	req, _ = http.NewRequest(
		"GET",
		server.URL + "/" + shortId,
		nil,
	)
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get originalUrl by shortUrl: %v", err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Waited status 202, but got %d", response.StatusCode)
	}

	bytes, _ = io.ReadAll(response.Body)
	defer response.Body.Close()

	sendedOriginalUrl := string(bytes)
	if sendedOriginalUrl != originalUrl {
		t.Errorf("Server bug when originalUrl saved incorrect")
	}
}
