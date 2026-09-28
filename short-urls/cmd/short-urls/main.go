package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"short-urls/internal/config"
	"short-urls/internal/handler"
	"short-urls/internal/service"

	"short-urls/internal/storage"
	"short-urls/internal/storage/memory"
	"short-urls/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		store storage.Storage
		err error
	)

	switch cfg.StorageType {
	case "memory":
		store = memory.NewMemoryStorage()
		log.Println("Local-memory storage starts")
	case "postgres":
		store, err = retryPostgresConnection(
			ctx,
			cfg.PostgresDSN,
			5,
			2 * time.Second,
		)
		if err != nil {
			log.Fatalf("Postgres storage error: %v", err)
		}
		log.Println("Postgres storage starts")
	default:
		log.Fatalf("StorageType is unknown: %s", cfg.StorageType)
	}
	defer store.Close()

	svc := service.NewService(store)
	hand := handler.NewHandler(svc, cfg.BaseURL)

	mux := http.NewServeMux()
	hand.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starts listen by addr = %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server starts shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5 * time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stops.")
}

func retryPostgresConnection(ctx context.Context, dsn string, attempts int, delay time.Duration) (storage.Storage, error) {
	var store storage.Storage
	var err error

	for i := 1; i <= attempts; i++ {
		store, err = postgres.NewPostgresStorage(ctx, dsn)
		if err == nil {
			break
		}
		log.Printf(
			"Postgres not connected (%d/%d) with error: %v",
			i, attempts, err,
		)
		time.Sleep(delay)
	}
	return store, err
}