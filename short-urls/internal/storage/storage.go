package storage

import (
	"errors"
	"context"
)

var (
	ErrUrlNotFound = errors.New("URL not found")
	ErrUrlExists = errors.New("URL already exists")
)

type Storage interface {
	Save(ctx context.Context, originalUrl, shortUrl string) error
	GetOriginal(ctx context.Context, shortUrl string) (string, error)
	GetShort(ctx context.Context, originalUrl string) (string, error)
	Close() error
}
