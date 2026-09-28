package service

import (
	"math/rand"
	"errors"
	"context"

	"short-urls/internal/storage"
)


const (
	Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	LenghtShortUrl = 10
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GenerateShortUrl(ctx context.Context, originalUrl string) (string, error) {
	shortUrl, err := s.storage.GetShort(ctx, originalUrl)
	if err == nil {
		return shortUrl, nil
	}
	if !errors.Is(err, storage.ErrUrlNotFound) {
		return "", err
	}

	for {
		shortUrl = generateRandomString(LenghtShortUrl)

		err = s.storage.Save(ctx, originalUrl, shortUrl)
		if err == nil {
			return shortUrl, nil
		}

		if errors.Is(err, storage.ErrUrlExists) {
			shortUrl, err := s.storage.GetShort(ctx, originalUrl)
			if err == nil {
				return shortUrl, nil
			}
			continue
		}
		return "", err
	}
}

func (s *Service) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	return s.storage.GetOriginal(ctx, shortUrl)
}

func generateRandomString(lenght int) string {
	bytes := make([]byte, lenght)
	for i := range bytes {
		bytes[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(bytes)
}
