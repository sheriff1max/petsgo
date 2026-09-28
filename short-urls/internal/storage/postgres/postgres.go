package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"short-urls/internal/storage"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, `
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

func (p *PostgresStorage) Save(ctx context.Context, originalUrl, shortUrl string) error {
	_, err := p.pool.Exec(
		ctx,
		`INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`,
		shortUrl,
		originalUrl,
	)

	if err != nil {
		if isUniqueError(err) {
			return storage.ErrUrlExists
		}
		return err
	}
	return nil
}

func (p *PostgresStorage) GetOriginal(ctx context.Context, shortUrl string) (string, error) {
	row := p.pool.QueryRow(
		ctx,
		`SELECT original_url FROM urls WHERE short_url = $1`,
		shortUrl,
	)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", err
	}
	return url, nil
}

func (p *PostgresStorage) GetShort(ctx context.Context, originalUrl string) (string, error) {
	row := p.pool.QueryRow(
		ctx,
		`SELECT short_url FROM urls WHERE original_url = $1`,
		originalUrl,
	)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrUrlNotFound
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
