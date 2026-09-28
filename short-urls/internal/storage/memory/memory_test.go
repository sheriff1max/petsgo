package memory

import (
	"errors"
	"testing"
	"context"

	"short-urls/internal/storage"
)

func TestMemorySaveAndGetOriginal(t *testing.T) {
	m := NewMemoryStorage()

	originalUrl := "https://example.com"
	shortUrl := "abc1234567"
	ctx := context.Background()
	
	err := m.Save(ctx, originalUrl, shortUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	orig, err := m.GetOriginal(ctx, shortUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig != originalUrl {
		t.Errorf("expected %s, got %s", originalUrl, orig)
	}
}

func TestMemorySaveAndGetByOriginal(t *testing.T) {
	m := NewMemoryStorage()

	originalUrl := "https://example.com"
	shortUrl := "abc1234567"
	ctx := context.Background()

	err := m.Save(ctx, originalUrl, shortUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	short, err := m.GetShort(ctx, originalUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if short != shortUrl {
		t.Errorf("expected %s, got %s", shortUrl, short)
	}
}

func TestMemoryGetOriginalNotFound(t *testing.T) {
	m := NewMemoryStorage()

	ctx := context.Background()
	_, err := m.GetOriginal(ctx, "nonexistent")
	if !errors.Is(err, storage.ErrUrlNotFound) {
		t.Errorf("expected ErrUrlNotFound, got %v", err)
	}
}

func TestMemory_GetByOriginal_NotFound(t *testing.T) {
	m := NewMemoryStorage()

	ctx := context.Background()
	_, err := m.GetShort(ctx, "https://nonexistent.com")
	if !errors.Is(err, storage.ErrUrlNotFound) {
		t.Errorf("expected ErrUrlNotFound, got %v", err)
	}
}

func TestMemoryDuplicateOriginal(t *testing.T) {
	m := NewMemoryStorage()

	originalUrl := "https://example.com"
	ctx := context.Background()

	err := m.Save(ctx, originalUrl, "abc1234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = m.Save(ctx, originalUrl, "xyz9876543")
	if !errors.Is(err, storage.ErrUrlExists) {
		t.Errorf("expected ErrUrlExists, got %v", err)
	}
}

func TestMemory_DuplicateShort(t *testing.T) {
	m := NewMemoryStorage()

	shortUrl := "abc1234567"
	ctx := context.Background()

	err := m.Save(ctx, "https://example1.com", shortUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = m.Save(ctx, "https://example2.com", shortUrl)
	if !errors.Is(err, storage.ErrUrlExists) {
		t.Errorf("expected ErrUrlExists, got %v", err)
	}
}
