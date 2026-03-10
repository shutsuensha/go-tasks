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

	"github.com/shutsuensha/go-tasks/internal/cache"
	"github.com/shutsuensha/go-tasks/internal/config"
	"github.com/shutsuensha/go-tasks/internal/db"
	"github.com/shutsuensha/go-tasks/internal/handler"
	"github.com/shutsuensha/go-tasks/internal/metrics"
	"github.com/shutsuensha/go-tasks/internal/middleware"
	"github.com/shutsuensha/go-tasks/internal/observability"
	"github.com/shutsuensha/go-tasks/internal/queue"
	"github.com/shutsuensha/go-tasks/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	"net/http/pprof"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	ctx := context.Background()

	shutdownTracing, err := observability.InitTracing(ctx, "go-tasks-api")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdownTracing(ctx)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	

	pgCfg, err := pgxpool.ParseConfig(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	pgCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		log.Fatal(err)
	}



	metrics.StartDBCollector(pool)

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("db not reachable:", err)
	}

	queries := db.New(pool)

	rdb := cache.NewRedisClient(cfg)

	rateLimiter := middleware.NewRedisRateLimiter(rdb, 10, time.Second)

	queueClient := queue.NewClient(cfg.RedisAddr)

	taskService := service.NewTaskService(pool, queries, rdb, queueClient)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Metrics)
	r.Use(rateLimiter.Middleware)

	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Post("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)

		r.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
		r.Get("/block", pprof.Handler("block").ServeHTTP)
		r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
		r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	})

	r.Handle("/metrics", promhttp.Handler())

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

	taskHandler := handler.NewTaskHandler(taskService)
	taskHandler.Register(r)

	server := &http.Server{
		Addr: ":" + cfg.HTTPPort,
		Handler: otelhttp.NewHandler(
			r,
			"http-server",
		),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		log.Println("server started on :" + cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

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