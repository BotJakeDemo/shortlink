// Package store holds short links and persists them to a JSON file.
package store

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
)

// Link is one short link.
type Link struct {
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

// Errors returned by Create.
var (
	ErrInvalidURL  = errors.New("url must be an absolute http or https URL")
	ErrInvalidSlug = errors.New("slug does not match the slug pattern")
	ErrSlugTaken   = errors.New("slug is already taken")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9-]{3,32}$`)

const (
	slugAlphabet  = "abcdefghijklmnopqrstuvwxyz0123456789"
	generatedLen  = 6
	maxGenerateTr = 100
)

// Store is a set of links backed by a JSON file.
type Store struct {
	mu    sync.Mutex
	path  string
	links map[string]Link
}

// Open loads the store at path. A missing file starts an empty store.
func Open(path string) (*Store, error) {
	s := &Store{path: path, links: map[string]Link{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var links []Link
	if err := json.Unmarshal(data, &links); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	for _, l := range links {
		s.links[l.Slug] = l
	}
	return s, nil
}

// Create validates and adds a link, generating a slug when slug is empty,
// and saves the store to its file.
func (s *Store) Create(rawURL, slug string) (Link, error) {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return Link{}, ErrInvalidURL
	}
	if slug != "" && !slugPattern.MatchString(slug) {
		return Link{}, ErrInvalidSlug
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if slug == "" {
		slug, err = s.generateSlug()
		if err != nil {
			return Link{}, err
		}
	} else if _, ok := s.links[slug]; ok {
		return Link{}, ErrSlugTaken
	}
	l := Link{Slug: slug, URL: rawURL}
	s.links[slug] = l
	if err := s.save(); err != nil {
		delete(s.links, slug)
		return Link{}, err
	}
	return l, nil
}

// Get returns the link with the given slug.
func (s *Store) Get(slug string) (Link, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.links[slug]
	return l, ok
}

// List returns every link, sorted by slug.
func (s *Store) List() []Link {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sorted()
}

func (s *Store) sorted() []Link {
	out := make([]Link, 0, len(s.links))
	for _, l := range s.links {
		out = append(out, l)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out
}

// generateSlug returns an unused random slug. The caller holds s.mu.
func (s *Store) generateSlug() (string, error) {
	n := big.NewInt(int64(len(slugAlphabet)))
	for try := 0; try < maxGenerateTr; try++ {
		b := make([]byte, generatedLen)
		for i := range b {
			r, err := rand.Int(rand.Reader, n)
			if err != nil {
				return "", err
			}
			b[i] = slugAlphabet[r.Int64()]
		}
		if _, ok := s.links[string(b)]; !ok {
			return string(b), nil
		}
	}
	return "", errors.New("could not generate an unused slug")
}

// save writes the whole file via a temp file and rename. The caller holds s.mu.
func (s *Store) save() error {
	data, err := json.Marshal(s.sorted())
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), filepath.Base(s.path)+".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), s.path)
}
