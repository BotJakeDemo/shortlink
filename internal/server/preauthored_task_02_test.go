//go:build preauthored

// Pre-authored evidence for task-02 of
// docs/shortlink/prds/prd-01-link-store-and-http-api.md. These assertions are
// fixed; the run implements the routes until they pass.
package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const seedLinks = `[{"slug":"zurdo","url":"https://github.com/ElOrlis/zurdo"},{"slug":"carol","url":"https://github.com/ElOrlis/carol"}]`

// newHandler returns a handler over a fresh data file holding seed, and the
// path of that file. An empty seed leaves the file missing.
func newHandler(t *testing.T, seed string) (http.Handler, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "links.json")
	if seed != "" {
		if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	h, err := New(path)
	if err != nil {
		t.Fatalf("New(%q): %v", path, err)
	}
	return h, path
}

func serve(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestServerRoutes(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		body         string
		wantStatus   int
		wantBody     string
		wantLocation string
	}{
		{name: "healthz", method: http.MethodGet, target: "/healthz", wantStatus: http.StatusOK, wantBody: "ok"},
		{name: "list", method: http.MethodGet, target: "/api/links", wantStatus: http.StatusOK, wantBody: `"slug":"zurdo"`},
		{name: "get", method: http.MethodGet, target: "/api/links/carol", wantStatus: http.StatusOK, wantBody: "https://github.com/ElOrlis/carol"},
		{name: "get unknown", method: http.MethodGet, target: "/api/links/nope", wantStatus: http.StatusNotFound},
		{name: "redirect", method: http.MethodGet, target: "/zurdo", wantStatus: http.StatusFound, wantLocation: "https://github.com/ElOrlis/zurdo"},
		{name: "redirect unknown", method: http.MethodGet, target: "/nope", wantStatus: http.StatusNotFound},
		{name: "create", method: http.MethodPost, target: "/api/links", body: `{"url":"https://example.com/new","slug":"newone"}`, wantStatus: http.StatusCreated, wantBody: `"slug":"newone"`},
		{name: "create generated slug", method: http.MethodPost, target: "/api/links", body: `{"url":"https://example.com/gen"}`, wantStatus: http.StatusCreated, wantBody: `"url":"https://example.com/gen"`},
		{name: "create not a url", method: http.MethodPost, target: "/api/links", body: `{"url":"not a url","slug":"badurl"}`, wantStatus: http.StatusBadRequest},
		{name: "create ftp url", method: http.MethodPost, target: "/api/links", body: `{"url":"ftp://example.com/f","slug":"ftpfile"}`, wantStatus: http.StatusBadRequest},
		{name: "create taken slug", method: http.MethodPost, target: "/api/links", body: `{"url":"https://example.com/x","slug":"zurdo"}`, wantStatus: http.StatusConflict},
		{name: "create bad slug", method: http.MethodPost, target: "/api/links", body: `{"url":"https://example.com/x","slug":"has space"}`, wantStatus: http.StatusBadRequest},
		{name: "create malformed body", method: http.MethodPost, target: "/api/links", body: `not json`, wantStatus: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := newHandler(t, seedLinks)
			rec := serve(h, tt.method, tt.target, tt.body)
			if rec.Code != tt.wantStatus {
				t.Fatalf("%s %s = %d, want %d (body %q)", tt.method, tt.target, rec.Code, tt.wantStatus, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body %q does not contain %q", rec.Body.String(), tt.wantBody)
			}
			if tt.wantStatus == http.StatusFound {
				if got := rec.Header().Get("Location"); got != tt.wantLocation {
					t.Errorf("Location = %q, want %q", got, tt.wantLocation)
				}
			}
			if !strings.HasPrefix(tt.target, "/api/") {
				return
			}
			if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if tt.wantStatus >= 400 {
				var e struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil || e.Error == "" {
					t.Errorf("error body %q is not {\"error\": \"<message>\"}", rec.Body.String())
				}
			}
		})
	}
}

func TestServerList(t *testing.T) {
	tests := []struct {
		name      string
		seed      string
		wantSlugs []string
	}{
		{"missing data file", "", []string{}},
		{"seeded", seedLinks, []string{"carol", "zurdo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := newHandler(t, tt.seed)
			rec := serve(h, http.MethodGet, "/api/links", "")
			if rec.Code != http.StatusOK {
				t.Fatalf("GET /api/links = %d, want %d", rec.Code, http.StatusOK)
			}
			body := strings.TrimSpace(rec.Body.String())
			if !strings.HasPrefix(body, "[") {
				t.Fatalf("body %q is not a JSON array", body)
			}
			var links []struct {
				Slug string `json:"slug"`
				URL  string `json:"url"`
			}
			if err := json.Unmarshal([]byte(body), &links); err != nil {
				t.Fatalf("body %q is not a JSON array of links: %v", body, err)
			}
			slugs := []string{}
			for _, l := range links {
				slugs = append(slugs, l.Slug)
			}
			if !reflect.DeepEqual(slugs, tt.wantSlugs) {
				t.Errorf("slugs = %v, want %v (sorted by slug)", slugs, tt.wantSlugs)
			}
		})
	}
}

func TestServerCreatePersists(t *testing.T) {
	h, path := newHandler(t, seedLinks)
	rec := serve(h, http.MethodPost, "/api/links", `{"url":"https://example.com/kept","slug":"kept"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/links = %d, want %d (body %q)", rec.Code, http.StatusCreated, rec.Body.String())
	}
	reopened, err := New(path)
	if err != nil {
		t.Fatalf("reopening %q: %v", path, err)
	}
	rec = serve(reopened, http.MethodGet, "/kept", "")
	if got := rec.Header().Get("Location"); rec.Code != http.StatusFound || got != "https://example.com/kept" {
		t.Fatalf("GET /kept after reopen = %d %q, want %d %q", rec.Code, got, http.StatusFound, "https://example.com/kept")
	}
}
