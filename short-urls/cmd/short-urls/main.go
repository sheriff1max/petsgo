package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"short-urls/internal/config"
	"short-urls/internal/handler"
	"short-urls/internal/service"
	"short-urls/internal/storage"
)

func main() {
	cfg := config.Load()

	var (
		store storage.Storage
		err error
	)

	switch cfg.StorageType {
	case "memory":
		store = storage.NewMemoryStorage()
		log.Println("Memory storage starts")
	case "postgres":
		store, err = storage.NewPostgresStorage(cfg.PostgresDSN)
		if err != nil {
			log.Fatalf("Postgres storage error: %v", err)
		}
		log.Println("Postgres storage starts")
	default:
		log.Fatalf("StorageType is unknown: %s", cfg.StorageType)
	}
	defer store.Close()

	svc := service.NewService(store)
	hand := handler.NewHandler(svc)

	mux := http.NewServeMux()
	hand.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	go func() {
		log.Println("Server starts listen by addr = %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server stops")
}
