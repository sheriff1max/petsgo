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

		countAtemps := 5
		for i := 1; i <= countAtemps; i++ {
			store, err = storage.NewPostgresStorage(cfg.PostgresDSN)
			if err == nil {
				break
			}
			log.Printf(
				"Postgres not connected (%d/%d) with error: %v",
				i, countAtemps, err,
			)
			time.Sleep(2 * time.Second)
		}
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

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stops.")
}
