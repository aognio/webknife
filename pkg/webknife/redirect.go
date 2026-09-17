package webknife

import (
	"fmt"
	"net/http"
	"net/url"
)

var validRedirectCodes = map[int]bool{
	301: true,
	302: true,
	303: true,
	307: true,
	308: true,
}

type RedirectConfig struct {
	StatusCode   int
	Target       string
	PreservePath bool
}

func NewRedirectHandler(cfg RedirectConfig) (http.Handler, error) {
	if !validRedirectCodes[cfg.StatusCode] {
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
