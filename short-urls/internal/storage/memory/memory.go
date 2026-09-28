package memory

import (
	"sync"
	"context"
	"short-urls/internal/storage"
)


type MemoryStorage struct {
	mu sync.RWMutex

	origToShort map[string]string
	shortToOrigin map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		origToShort: make(map[string]string),
		shortToOrigin: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(ctx context.Context, originalUrl, shortUrl string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.origToShort[originalUrl]; ok {
		return storage.ErrUrlExists
	}

	if _, ok := m.shortToOrigin[shortUrl]; ok {
		return storage.ErrUrlExists
	}

	m.origToShort[originalUrl] = shortUrl
	m.shortToOrigin[shortUrl] = originalUrl
	return nil
}

func (m *MemoryStorage) GetOriginal(ctx context.Context, shortUrl string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	url, ok := m.shortToOrigin[shortUrl]
	if !ok {
		return "", storage.ErrUrlNotFound
	}
	return url, nil
}

func (m *MemoryStorage) GetShort(ctx context.Context, originalUrl string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	url, ok := m.origToShort[originalUrl]
	if !ok {
		return "", storage.ErrUrlNotFound
	}
	return url, nil
}

func (m *MemoryStorage) Close() error {
	return nil
}
