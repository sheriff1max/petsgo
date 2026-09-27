package storage

import "sync"


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

func (m *MemoryStorage) Save(originalUrl, shortUrl string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.origToShort[originalUrl]; ok {
		return ErrUrlExists
	}

	if _, ok := m.shortToOrigin[shortUrl]; ok {
		return ErrUrlExists
	}

	m.origToShort[originalUrl] = shortUrl
	m.shortToOrigin[shortUrl] = originalUrl
	return nil
}

func (m *MemoryStorage) GetOriginal(shortUrl string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	url, ok := m.shortToOrigin[shortUrl]
	if !ok {
		return "", ErrUrlNotFound
	}
	return url, nil
}

func (m *MemoryStorage) GetShort(originalUrl string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	url, ok := m.origToShort[originalUrl]
	if !ok {
		return "", ErrUrlNotFound
	}
	return url, nil
}

func (m *MemoryStorage) Close() error {
	return nil
}
