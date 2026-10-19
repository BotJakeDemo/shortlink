// Pre-authored evidence for task-01 of
// docs/shortlink/prds/prd-01-link-store-and-http-api.md. These assertions are
// fixed; the run implements the store until they pass.
package store

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestStoreCreate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		slug    string
		wantErr error
	}{
		{"explicit slug", "https://github.com/ElOrlis/zurdo", "zurdo", nil},
		{"http url", "http://example.com/a?b=c", "plainhttp", nil},
		{"generated slug", "https://example.com/", "", nil},
		{"hyphenated slug", "https://example.com/h", "my-link", nil},
		{"taken slug", "https://example.com/other", "taken", ErrSlugTaken},
		{"relative url", "/just/a/path", "relative", ErrInvalidURL},
		{"ftp url", "ftp://example.com/file", "ftpfile", ErrInvalidURL},
		{"empty url", "", "emptyurl", ErrInvalidURL},
		{"slug with space", "https://example.com/", "has space", ErrInvalidSlug},
		{"slug with slash", "https://example.com/", "a/b", ErrInvalidSlug},
		{"slug too short", "https://example.com/", "ab", ErrInvalidSlug},
		{"slug too long", "https://example.com/", strings.Repeat("a", 33), ErrInvalidSlug},
		{"uppercase slug", "https://example.com/", "Zurdo", ErrInvalidSlug},
	}
	generated := regexp.MustCompile(`^[a-z0-9]{6}$`)
	s, err := Open(filepath.Join(t.TempDir(), "links.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create("https://example.com/first", "taken"); err != nil {
		t.Fatalf("seeding slug %q: %v", "taken", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(tt.url, tt.slug)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create(%q, %q) error = %v, want %v", tt.url, tt.slug, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.URL != tt.url {
				t.Errorf("URL = %q, want %q", got.URL, tt.url)
			}
			if tt.slug != "" && got.Slug != tt.slug {
				t.Errorf("Slug = %q, want %q", got.Slug, tt.slug)
			}
			if tt.slug == "" && !generated.MatchString(got.Slug) {
				t.Errorf("generated Slug = %q, want 6 lowercase letters/digits", got.Slug)
			}
			if stored, ok := s.Get(got.Slug); !ok || stored != got {
				t.Errorf("Get(%q) = %v, %v; want %v, true", got.Slug, stored, ok, got)
			}
		})
	}
	// taken + the four successful creates; rejected creates must not be stored.
	if got := len(s.List()); got != 5 {
		t.Errorf("len(List()) = %d, want 5", got)
	}
}

func TestStoreOpenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "links.json")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open(%q) on a missing file: %v", path, err)
	}
	if got := s.List(); len(got) != 0 {
		t.Fatalf("List() on a fresh store = %v, want empty", got)
	}
}

func TestStoreOpenUnreadable(t *testing.T) {
	dir := t.TempDir()
	garbage := filepath.Join(dir, "garbage.json")
	if err := os.WriteFile(garbage, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		path string
	}{
		{"directory", dir},
		{"malformed json", garbage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Open(tt.path); err == nil {
				t.Fatalf("Open(%q) = nil error, want an error", tt.path)
			}
		})
	}
}

func TestStoreReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "links.json")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range []Link{
		{Slug: "zurdo", URL: "https://github.com/ElOrlis/zurdo"},
		{Slug: "carol", URL: "https://github.com/ElOrlis/carol"},
		{Slug: "bravo", URL: "http://example.com/b?q=1&r=2"},
	} {
		if _, err := s.Create(l.URL, l.Slug); err != nil {
			t.Fatalf("Create(%q, %q): %v", l.URL, l.Slug, err)
		}
	}
	want := []Link{
		{Slug: "bravo", URL: "http://example.com/b?q=1&r=2"},
		{Slug: "carol", URL: "https://github.com/ElOrlis/carol"},
		{Slug: "zurdo", URL: "https://github.com/ElOrlis/zurdo"},
	}
	if got := s.List(); !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v (sorted by slug)", got, want)
	}
	reloaded, err := Open(path)
	if err != nil {
		t.Fatalf("reopening %q: %v", path, err)
	}
	if got := reloaded.List(); !reflect.DeepEqual(got, want) {
		t.Fatalf("after reload List() = %v, want %v", got, want)
	}
}
