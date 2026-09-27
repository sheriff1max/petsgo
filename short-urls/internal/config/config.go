package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr string
	StorageType string
	PostgresDSN string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(
		&cfg.Addr,
		"a",
		getEnv("SERVER_ADDR", ":8080"),
		"HTTP server addres",
	)
	flag.StringVar(
		&cfg.StorageType,
		"s",
		getEnv("STORAGE_TYPE", "memory"),
		"Storage type: memory or postgres",
	)
	flag.StringVar(
		&cfg.PostgresDSN,
		"d",
		getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/shortener?sslmode=disable"),
		"PostgreSQL DSN",
	)

	flag.Parse()
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}