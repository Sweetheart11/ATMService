package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/exp/slog"

	"github.com/Sweetheart11/ATMService/internal/config"
	"github.com/Sweetheart11/ATMService/internal/http-server/handlers/urls/balance"
	"github.com/Sweetheart11/ATMService/internal/http-server/handlers/urls/create"
	"github.com/Sweetheart11/ATMService/internal/http-server/handlers/urls/deposit"
	"github.com/Sweetheart11/ATMService/internal/http-server/handlers/urls/withdraw"
	"github.com/Sweetheart11/ATMService/internal/storage/sliceStorage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting ATM service",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	storage, err := sliceStorage.New()
	// error handling for a more complex implementation
	if err != nil {
		log.Error("failed initializing storage", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	router.Post("/accounts", create.New(log, &storage))
	router.Post("/accounts/{id}/deposit", deposit.New(log, &storage))
	router.Post("/accounts/{id}/withdraw", withdraw.New(log, &storage))
	router.Get("/accounts/{id}/balance", balance.New(log, &storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.StringValue(err.Error()))

		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
