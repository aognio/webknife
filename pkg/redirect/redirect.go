// Package redirect provides a handler that redirects requests to a target URL.
//
// This package has no dependencies on other Webknife feature packages.
package redirect

import (
	"fmt"
	"net/http"
	"net/url"
)

var validCodes = map[int]bool{301: true, 302: true, 303: true, 307: true, 308: true}

// Config holds configuration for the redirect handler.
type Config struct {
	StatusCode   int
	Target       string
	PreservePath bool
}

// New returns an http.Handler that redirects to the configured target.
func New(cfg Config) (http.Handler, error) {
	if !validCodes[cfg.StatusCode] {
		return nil, fmt.Errorf("invalid redirect status code: %d (must be 301, 302, 303, 307, or 308)", cfg.StatusCode)
	}
	targetURL, err := url.Parse(cfg.Target)
	if err != nil {
		return nil, fmt.Errorf("invalid redirect target: %w", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var redirectURL *url.URL
		if cfg.PreservePath {
			redirectURL = &url.URL{
				Scheme:   targetURL.Scheme,
				Host:     targetURL.Host,
				Path:     r.URL.Path,
				RawQuery: r.URL.RawQuery,
				Fragment: r.URL.Fragment,
			}
		} else {
			redirectURL = targetURL
		}
		http.Redirect(w, r, redirectURL.String(), cfg.StatusCode)
	}), nil
}
