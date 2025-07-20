package server

import (
	"encoding/json"
	"errors"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/models"
	"net/http"
	"strconv"
)

type createReqPayload struct {
	Url string `json:"url"`
}
type createRespPayload struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type feedResp struct {
	Id  int    `json:"id,omitempty"`
	Url string `json:"url"`
}
type feedNewsResp struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

func (s server) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqPayload createReqPayload
	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	if reqPayload.Url == "" {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	// check that the URL is not added before
	exist, err := s.feedRepository.GetByUrl(r.Context(), reqPayload.Url)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			s.serverError(w, r, err)
			return
		}
	}
	if exist.ID != 0 {
		s.logger.Info("The URL already exists.", "url", reqPayload.Url)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	feed := models.Feed{
		Url: reqPayload.Url,
	}

	id, err := s.feedRepository.Create(ctx, feed)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	s.worker.DownloadInBackground(id, reqPayload.Url)

	resp := createRespPayload{
		Id:  id,
		Url: feed.Url,
	}
	s.responseJSON(w, resp)
}

func (s server) deleteFeedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	err = s.feedRepository.Delete(r.Context(), id)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
}

func (s server) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.feedRepository.GetAll(r.Context())
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	var resp = make([]feedResp, len(feeds))
	for i, feed := range feeds {
		resp[i] = feedResp{
			Id:  feed.ID,
			Url: feed.Url,
		}
	}
	s.responseJSON(w, resp)
}

func (s server) getFeedByIdHandler(w http.ResponseWriter, r *http.Request) {
	feedId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	feed, err := s.feedRepository.GetById(r.Context(), feedId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			s.serverError(w, r, err)
		}
		return
	}

	resp := feedResp{
		Url: feed.Url,
	}
	s.responseJSON(w, resp)
}

func (s server) getFeedNewsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	items, err := s.itemRepository.GetAllInFeed(r.Context(), id)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	var resp = make([]feedNewsResp, len(items))
	for i, item := range items {
		resp[i] = feedNewsResp{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		}
	}
	s.responseJSON(w, resp)
}
