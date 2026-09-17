// Package static serves files from a directory.
//
// It uses Go's standard http.FileServer with path traversal protection.
// This package has no dependencies on other Webknife packages.
package static

import (
	"net/http"
	"path/filepath"
)

// Config holds configuration for the static file server.
type Config struct {
	// Root is the directory to serve files from.
	Root string
}

// New returns an http.Handler that serves static files from cfg.Root.
func New(cfg Config) (http.Handler, error) {
	absRoot, err := filepath.Abs(cfg.Root)
	if err != nil {
		return nil, err
	}
	return http.StripPrefix("/", http.FileServer(http.Dir(absRoot))), nil
}
