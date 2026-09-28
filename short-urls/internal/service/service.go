package service

import (
	"math/rand"
	"errors"

	"short-urls/internal/storage"
)


const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	LenghtShortUrl = 10
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GenerateShortUrl(originalUrl string) (string, error) {
	shortUrl, err := s.storage.GetShort(originalUrl)
	if err == nil {
		return shortUrl, nil
	}
	if !errors.Is(err, storage.ErrUrlNotFound) {
		return "", err
	}

	for {
		shortUrl = generateRandomString(LenghtShortUrl)

		err = s.storage.Save(originalUrl, shortUrl)
		if err == nil {
			return shortUrl, nil
		}

		if errors.Is(err, storage.ErrUrlExists) {
			shortUrl, err := s.storage.GetShort(originalUrl)
			if err == nil {
				return shortUrl, nil
			}
			continue
		}
		return "", err
	}
}

func (s *Service) GetOriginalUrl(shortUrl string) (string, error) {
	return s.storage.GetOriginal(shortUrl)
}

func generateRandomString(lenght int) string {
	bytes := make([]byte, lenght)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}
