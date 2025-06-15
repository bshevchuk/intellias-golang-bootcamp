package main

import (
	"context"
	"fmt"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/database"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/downloader"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/parser"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

// const defaultRssUrl = "https://dou.ua/feed/"
const defaultRssUrl = "https://news.ycombinator.com/rss"

const defaultDatabaseUrl = "postgres://pguser:pgpassword@localhost:5432/pgdb?sslmode=disable"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create database
	dbpool, err := pgxpool.New(ctx, defaultDatabaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	db := database.NewPostgres(dbpool)
	log.Printf("Connecting to database %s", defaultDatabaseUrl)

	// Delete all data in database to have "fresh" database (for development purpose only)
	err = db.DeleteAll(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to delete all from db: %v\n", err)
		os.Exit(2)
	}
	log.Println("Deleted all data in database")

	// Download content
	content, err := downloader.Download(defaultRssUrl)
	if err != nil {
		fmt.Printf("exit with error: %v", err)
		os.Exit(1)
	}
	log.Printf("Downloaded %d bytes from %s", len(content), defaultRssUrl)

	// Parse content as RSS
	rss, err := parser.ParseRss(content)
	if err != nil {
		fmt.Printf("exit with error: %v", err)
		os.Exit(1)
	}
	log.Println("Parsed RSS")

	// Show RSS and save into database
	fmt.Printf("%s\n", rss.Channel.Title)
	for _, item := range rss.Channel.Items {
		err := db.CreateItem(ctx, item.Title, item.Link, item.Description)
		if err != nil {
			fmt.Printf("error when creating item %s error: %v", item.Title, err)
		}

		fmt.Printf("\t %s\n", item.Title)
		fmt.Printf("\t %s\n", item.Link)
		fmt.Println()
	}

	log.Println("Done")
}
