// Package server serves the shortlink HTTP API.
package server

import "net/http"

// New returns the handler for every shortlink route, backed by the JSON
// store at dataPath. A missing file starts an empty store; an unreadable one
// is an error.
func New(dataPath string) (http.Handler, error) { panic("todo: task-02") }
