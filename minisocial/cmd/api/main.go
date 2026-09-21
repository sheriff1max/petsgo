// @title MiniSocial API
// @version 1.0
// @description Мини соц-сеть на Go
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"minisocial/internal/config"
	"minisocial/internal/handler"
	"minisocial/internal/repository"
	"minisocial/internal/service"

	_ "github.com/swaggo/files/v2"
    httpSwagger "github.com/swaggo/http-swagger/v2"
    _ "minisocial/docs"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, repository.GetDSN(cfg.DB))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	repo := repository.NewPostgresRepository(dbpool)
	svc := service.NewSocialService(repo)
	h := handler.NewHTTPHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8080/swagger/doc.json"), 
    ))

	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
