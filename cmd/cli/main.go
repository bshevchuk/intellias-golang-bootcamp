package main

import (
	"context"
	"fmt"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/downloader"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/models"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/parser"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
)

// const defaultRssUrl = "https://dou.ua/feed/"
const defaultRssUrl = "https://news.ycombinator.com/rss"

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

	itemRepository := repository.NewItemRepository(dbpool)
	logger.Info("Connecting to database", slog.Any("databaseUrl", defaultDatabaseUrl))

	// Delete all data in database to have "fresh" database (for development purpose only)
	err = itemRepository.DeleteAll(ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(2)
	}
	logger.Info("Deleted all data in database")

	// Download content
	content, err := downloader.Download(defaultRssUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Debug("Downloaded ", slog.Any("bytes", len(content)), slog.Any("url", defaultRssUrl))

	// Parse content as RSS
	rss, err := parser.ParseRss(content)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Parsed RSS")

	// Show RSS and save into database
	fmt.Printf("%s\n", rss.Channel.Title)
	for _, item := range rss.Channel.Items {
		err := itemRepository.Create(ctx, models.Item{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		})
		if err != nil {
			fmt.Printf("error when creating item %s error: %v", item.Title, err)
		}

		fmt.Printf("\t %s\n", item.Title)
		fmt.Printf("\t %s\n", item.Link)
		fmt.Println()
	}

	logger.Info("Done")
}
