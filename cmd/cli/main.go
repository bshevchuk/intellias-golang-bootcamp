package main

import (
	"fmt"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/downloader"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/parser"
	"os"
)

const defaultRssUrl = "https://dou.ua/feed/"

func main() {

	// Download content
	content, err := downloader.Download(defaultRssUrl)
	if err != nil {
		fmt.Printf("exit with error: %v", err)
		os.Exit(1)
	}

	// Parse as RSS
	rss, err := parser.ParseRss(content)
	if err != nil {
		fmt.Printf("exit with error: %v", err)
		os.Exit(1)
	}

	// Show RSS
	fmt.Printf("%s\n", rss.Channel.Title)
	for _, item := range rss.Channel.Items {
		fmt.Printf("\t %s\n"+
			"\t %s\n\n", item.Title, item.Link)
	}
}
