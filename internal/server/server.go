// Package server serves the shortlink HTTP API.
package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BotJakeDemo/shortlink/internal/store"
)

// New returns the handler for every shortlink route, backed by the JSON
// store at dataPath. A missing file starts an empty store; an unreadable one
// is an error.
func New(dataPath string) (http.Handler, error) {
	s, err := store.Open(dataPath)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /api/links", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL  string `json:"url"`
			Slug string `json:"slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}
		l, err := s.Create(req.URL, req.Slug)
		switch {
		case errors.Is(err, store.ErrInvalidURL), errors.Is(err, store.ErrInvalidSlug):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, store.ErrSlugTaken):
			writeError(w, http.StatusConflict, err.Error())
		case err != nil:
			writeError(w, http.StatusInternalServerError, err.Error())
		default:
			writeJSON(w, http.StatusCreated, l)
		}
	})
	mux.HandleFunc("GET /api/links", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.List())
	})
	mux.HandleFunc("GET /api/links/{slug}", func(w http.ResponseWriter, r *http.Request) {
		l, ok := s.Get(r.PathValue("slug"))
		if !ok {
			writeError(w, http.StatusNotFound, "link not found")
			return
		}
		writeJSON(w, http.StatusOK, l)
	})
	mux.HandleFunc("GET /{slug}", func(w http.ResponseWriter, r *http.Request) {
		l, ok := s.Get(r.PathValue("slug"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, l.URL, http.StatusFound)
	})
	return mux, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
