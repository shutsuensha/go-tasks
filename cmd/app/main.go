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
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shutsuensha/go-tasks/internal/config"
	"github.com/shutsuensha/go-tasks/internal/db"
	"github.com/shutsuensha/go-tasks/internal/handler"
	"github.com/shutsuensha/go-tasks/internal/middleware"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()


	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("db not reachable:", err)
	}

	queries := db.New(pool)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recovery(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("db not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	taskHandler := handler.NewTaskHandler(queries)
	taskHandler.Register(r)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	// запуск сервера
	go func() {
		log.Println("server started on :" + cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error:", err)
	}

	log.Println("server stopped")
}