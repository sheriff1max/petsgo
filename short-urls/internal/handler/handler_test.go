package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"short-urls/internal/service"
	"short-urls/internal/storage/memory"
)


func setupTestServer() (*httptest.Server, string) {
	BaseURL := "http://localhost:8080"

	store := memory.NewMemoryStorage()
	svc := service.NewService(store)
	h := NewHandler(svc, BaseURL)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return httptest.NewServer(mux), BaseURL
}

func TestHandlerGenerateShortUrlSuccess(t *testing.T) {
	server, BaseURL := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader("https://github.com"),
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.HasPrefix(string(body), BaseURL) {
		t.Errorf("expected response to start with base URL, got %s", string(body))
	}
}

func TestHandlerGenerateShortUrlEmptyBody(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(""),
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandlerGenerateShortUrlWhitespaceOnly(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader("   \n\t  "),
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandlerGenerateShortUrlIdempotent(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	url := "https://example.com/idempotent"

	// Первый запрос
	req1, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(url),
	)
	resp1, _ := http.DefaultClient.Do(req1)
	body1, _ := io.ReadAll(resp1.Body)
	defer resp1.Body.Close()

	// Второй запрос с тем же url
	req2, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(url),
	)
	resp2, _ := http.DefaultClient.Do(req2)
	body2, _ := io.ReadAll(resp2.Body)
	defer resp2.Body.Close()

	if string(body1) != string(body2) {
		t.Errorf("expected same response for same URL, got %s and %s", string(body1), string(body2))
	}
}

func TestHandlerGetOriginalUrlNotFound(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest(
		"GET",
		server.URL + "/nonexistent",
		nil,
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestResolve_EmptyShort(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL + "/", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestFullPipeline(t *testing.T) {
	server, _ := setupTestServer()
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
