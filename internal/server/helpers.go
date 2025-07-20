package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (s server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	s.logger.Error(err.Error(), slog.Any("method", method), slog.Any("uri", uri), slog.Any("trace", trace))

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s server) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (s server) responseJSON(w http.ResponseWriter, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		s.logger.Error(err.Error())
	}
}
