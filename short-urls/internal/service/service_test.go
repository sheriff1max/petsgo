package service

import (
	"strings"
	"testing"
	"context"

	"short-urls/internal/storage/memory"
)

func TestGenerateRandomStringLength(t *testing.T) {
	s := generateRandomString(LenghtShortUrl)
	if len(s) != LenghtShortUrl {
		t.Errorf("expected length %d, got %d", LenghtShortUrl, len(s))
	}
}

func TestGenerateRandomStringCharset(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := generateRandomString(LenghtShortUrl)
		for _, ch := range s {
			if !strings.ContainsRune(Charset, ch) {
				t.Errorf("character '%c' is not in allowed alphabet", ch)
			}
		}
	}
}

func TestGenerateShortUrl(t *testing.T) {
	store := memory.NewMemoryStorage()
	svc := NewService(store)

	short, err := svc.GenerateShortUrl(
		context.Background(),
		"https://example.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(short) != LenghtShortUrl {
		t.Errorf("expected short URL length 10, got %d", len(short))
	}
}

func TestGenerateTwiceShortUrl(t *testing.T) {
	store := memory.NewMemoryStorage()
	svc := NewService(store)

	url := "https://example.com/same"

	short1, err := svc.GenerateShortUrl(context.Background(), url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	short2, err := svc.GenerateShortUrl(context.Background(), url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if short1 != short2 {
		t.Errorf("expected same short URL for same original, got %s and %s", short1, short2)
	}
}

func TestGenerateShortUrlForDifferentOriginalUrls(t *testing.T) {
	store := memory.NewMemoryStorage()
	svc := NewService(store)

	short1, _ := svc.GenerateShortUrl(
		context.Background(),
		"https://example.com/one",
	)
	short2, _ := svc.GenerateShortUrl(
		context.Background(),
		"https://example.com/two",
	)

	if short1 == short2 {
		t.Errorf("different URLs should have different short links, both got %s", short1)
	}
}

func TestGetOriginalUrl(t *testing.T) {
	store := memory.NewMemoryStorage()
	svc := NewService(store)

	original := "https://example.com/resolve"
	short, _ := svc.GenerateShortUrl(context.Background(), original)
	gettedOriginal, err := svc.GetOriginalUrl(context.Background(), short)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gettedOriginal != original {
		t.Errorf("expected %s, got %s", original, gettedOriginal)
	}
}

func TestGetOriginalUrlUnknown(t *testing.T) {
	store := memory.NewMemoryStorage()
	svc := NewService(store)

	_, err := svc.GetOriginalUrl(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent short URL, got nil")
	}
}
