package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			short_url VARCHAR(10) PRIMARY KEY,
			original_url TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{pool: pool}, nil
}

func (p *PostgresStorage) Save(originalUrl, shortUrl string) error {
	_, err := p.pool.Exec(
		context.Background(),
		`INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`,
		shortUrl,
		originalUrl,
	)

	if err != nil {
		if isUniqueError(err) {
			return ErrUrlExists
		}
		return err
	}
	return nil
}

func (p *PostgresStorage) GetOriginal(shortUrl string) (string, error) {
	row := p.pool.QueryRow(
		context.Background(),
		`SELECT original_url FROM urls WHERE short_url = $1`,
		shortUrl,
	)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUrlNotFound
		}
		return "", err
	}
	return url, nil
}

func (p *PostgresStorage) GetShort(originalUrl string) (string, error) {
	row := p.pool.QueryRow(
		context.Background(),
		`SELECT short_url FROM urls WHERE original_url = $1`,
		originalUrl,
	)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUrlNotFound
		}
		return "", err
	}
	return url, nil
}

func (p *PostgresStorage) Close() error {
	p.pool.Close()
	return nil
}

func isUniqueError(err error) bool {
	// проверка кода ошибки postgresql 23505
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
