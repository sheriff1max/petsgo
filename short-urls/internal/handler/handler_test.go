package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"encoding/json"

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

	request := RequestOriginalUrl{OriginalUrl: "https://github.com"}
	requestBody, _ := json.Marshal(request)

	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(string(requestBody)),
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var responseJson ResponseShortUrl
	if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	if !strings.HasPrefix(responseJson.ShortUrl, BaseURL) {
		t.Errorf(
			"expected response to start with base URL, got %s",
			responseJson.ShortUrl,
		)
	}
}

func TestHandlerGenerateShortUrlEmptyBody(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	request := RequestOriginalUrl{OriginalUrl: ""}
	requestBody, _ := json.Marshal(request)

	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(string(requestBody)),
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

	request := RequestOriginalUrl{OriginalUrl: "   \n\t  "}
	requestBody, _ := json.Marshal(request)

	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(string(requestBody)),
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

	request := RequestOriginalUrl{OriginalUrl: url}
	requestBody, _ := json.Marshal(request)

	// Первый запрос
	req1, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(string(requestBody)),
	)
	resp1, _ := http.DefaultClient.Do(req1)
	var response1 ResponseShortUrl
	if err := json.NewDecoder(resp1.Body).Decode(&response1); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	defer resp1.Body.Close()

	// Второй запрос с тем же url
	req2, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(string(requestBody)),
	)
	resp2, _ := http.DefaultClient.Do(req2)
	var response2 ResponseShortUrl
	if err := json.NewDecoder(resp2.Body).Decode(&response2); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	defer resp2.Body.Close()

	if response1.ShortUrl != response2.ShortUrl {
		t.Errorf("expected same response for same URL, got %s and %s", response1.ShortUrl, response2.ShortUrl)
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

	request := RequestOriginalUrl{OriginalUrl: "https://github.com"}
	requestBody, _ := json.Marshal(request)

	// post метод
	req, _ := http.NewRequest(
		"POST",
		server.URL + "/generate",
		strings.NewReader(string(requestBody)),
	)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed generation to shortUrl: %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Waited status 201, but gor %d", response.StatusCode)
	}

	var responseJson ResponseShortUrl
	if err := json.NewDecoder(response.Body).Decode(&responseJson); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	defer response.Body.Close()

	shortUrl := responseJson.ShortUrl
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

	var responseJson2 ResponseOriginalUrl
	if err := json.NewDecoder(response.Body).Decode(&responseJson2); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	defer response.Body.Close()

	sendedOriginalUrl := responseJson2.OriginalUrl
	if sendedOriginalUrl != request.OriginalUrl {
		t.Errorf("Server bug when originalUrl saved incorrect")
	}
}
