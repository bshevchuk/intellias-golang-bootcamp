package server

import (
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/repository"
	"log/slog"
	"net/http"
)

type server struct {
	itemRepository repository.ItemRepository
	feedRepository repository.FeedRepository
	logger         *slog.Logger
	worker         backgroundWorker
}

type backgroundWorker interface {
	DownloadInBackground(feedId int, feedUrl string)
}

func NewServer(itemRepository repository.ItemRepository, feedRepository repository.FeedRepository, logger *slog.Logger, bw backgroundWorker) *server {
	return &server{
		itemRepository: itemRepository,
		feedRepository: feedRepository,
		logger:         logger,
		worker:         bw,
	}
}

func (s server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /feed", s.createFeedHandler)
	mux.HandleFunc("DELETE /feed/{id}", s.deleteFeedHandler)
	mux.HandleFunc("GET /feed", s.getFeedsHandler)
	mux.HandleFunc("GET /feed/{id}", s.getFeedByIdHandler)
	mux.HandleFunc("GET /feed/{id}/news", s.getFeedNewsHandler)

	return s.recoverPanic(
		s.logRequest(
			mux,
		),
	)
}
