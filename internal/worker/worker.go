package worker

import (
	"context"
	"fmt"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/downloader"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/models"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/parser"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/repository"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
)

type downloadJob struct {
	feedId int
	url    string
}
type worker struct {
	itemRepository repository.ItemRepository
	feedRepository repository.FeedRepository
	ctx            context.Context
	logger         *slog.Logger
	downloadCh     chan downloadJob
}

func NewWorker(ctx context.Context, itemRepository repository.ItemRepository, feedRepository repository.FeedRepository, logger *slog.Logger) *worker {
	return &worker{
		ctx:            ctx,
		itemRepository: itemRepository,
		feedRepository: feedRepository,
		logger:         logger,
		downloadCh:     make(chan downloadJob),
	}
}

func (w *worker) DownloadInBackground(feedId int, feedUrl string) {
	w.logger.Info("DownloadInBackground added",
		slog.Any("feedId", feedId),
		slog.Any("url", feedUrl),
		slog.Any("goid", goid()))
	w.downloadCh <- downloadJob{
		feedId: feedId,
		url:    feedUrl,
	}
}

func (w *worker) Start() {
	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("Worker stopped", slog.Any("goid", goid()))
			return
		case job := <-w.downloadCh:
			w.downloadNewFeed(job.feedId, job.url)
		}
	}
}

func (w *worker) downloadNewFeed(feedId int, feedUrl string) {
	w.logger.Info("Downloading new feed. Start",
		slog.Any("feedId", feedId),
		slog.Any("url", feedUrl),
		slog.Any("goid", goid()))
	content, err := downloader.Download(feedUrl)
	if err != nil {
		w.logger.Error(err.Error())
	}
	w.logger.Info("Downloaded ",
		slog.Any("bytes", len(content)),
		slog.Any("url", feedUrl),
		slog.Any("goid", goid()))

	rss, err := parser.ParseRss(content)
	if err != nil {
		w.logger.Error(err.Error())
	}

	for _, item := range rss.Channel.Items {
		err := w.itemRepository.Create(w.ctx, models.Item{
			FeedID:      feedId,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		})
		if err != nil {
			w.logger.Error(err.Error(),
				slog.Any("goid", goid()))
		}
	}
	w.logger.Info("Downloading new feed.. DONE",
		slog.Any("goid", goid()))
}

// goid is used for debugging goroutines. this SHOULD NOT be used in production code
//
//	https://gist.github.com/metafeather/3615b23097836bc36579100dac376906
func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
