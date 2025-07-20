package main

import (
	"context"
	"fmt"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/repository"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const defaultDatabaseUrl = "postgres://pguser:pgpassword@localhost:5432/pgdb?sslmode=disable"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Create database
	dbpool, err := pgxpool.New(ctx, defaultDatabaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	log.Printf("Connecting to database %s", defaultDatabaseUrl)

	itemRepository := repository.NewItemRepository(dbpool)
	feedRepository := repository.NewFeedRepository(dbpool)
	s := server.NewServer(itemRepository, feedRepository, logger)

	port := 3000
	addr := fmt.Sprintf(":%d", port)
	log.Printf("starting server on %s", addr)

	err = http.ListenAndServe(addr, s.Routes())
	if err != nil {
		log.Fatal(err)
		return
	}
}
