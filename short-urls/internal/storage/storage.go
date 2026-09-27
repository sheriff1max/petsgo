package storage

import "errors"

var (
	ErrUrlNotFound = errors.New("URL not found")
	ErrUrlExists = errors.New("URL already exists")
)

type Storage interface {
	Save(originalUrl, shortUrl string) error
	GetOriginal(shortUrl string) (string, error)
	GetShort(originalUrl string) (string, error)
	Close() error
}
