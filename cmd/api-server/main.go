package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/repository"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/server"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultDSN = "postgres://pguser:pgpassword@localhost:5432/pgdb?sslmode=disable"
const defaultPort = 3000

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	// Create database pool
	dsn := defaultDSN
	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	logger.Info("Connecting to database", "dsn", dsn)

	itemRepository := repository.NewItemRepository(dbpool)
	feedRepository := repository.NewFeedRepository(dbpool)
	w := worker.NewWorker(ctx, itemRepository, feedRepository, logger)
	s := server.NewServer(itemRepository, feedRepository, logger, w)

	port := defaultPort
	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting server", "addr", addr)

	srv := http.Server{
		Addr:    addr,
		Handler: s.Routes(),
	}

	// start server
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Info("Stopped serving new connections")
	}()

	// start background worker
	go w.Start()

	// === Graceful shutdown section ===
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Graceful shutdown complete")
}
