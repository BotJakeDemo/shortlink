// Package store holds short links and persists them to a JSON file.
package store

import "errors"

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

// Store is a set of links backed by a JSON file.
type Store struct{}

// Open loads the store at path. A missing file starts an empty store.
func Open(path string) (*Store, error) { panic("todo: task-01") }

// Create validates and adds a link, generating a slug when slug is empty,
// and saves the store to its file.
func (s *Store) Create(url, slug string) (Link, error) { panic("todo: task-01") }

// Get returns the link with the given slug.
func (s *Store) Get(slug string) (Link, bool) { panic("todo: task-01") }

// List returns every link, sorted by slug.
func (s *Store) List() []Link { panic("todo: task-01") }
