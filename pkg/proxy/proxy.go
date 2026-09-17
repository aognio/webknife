// Package proxy provides a reverse proxy handler using Go's httputil.ReverseProxy.
//
// This package has no dependencies on other Webknife feature packages.
package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Config holds configuration for the reverse proxy.
type Config struct {
	// Upstream is the target URL to proxy requests to.
	Upstream string
}

// New returns an http.Handler that reverse-proxies requests to the upstream.
func New(cfg Config) (http.Handler, error) {
	upstreamURL, err := url.Parse(cfg.Upstream)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = upstreamURL.Scheme
		req.URL.Host = upstreamURL.Host
		req.Host = upstreamURL.Host
	}
	return proxy, nil
}
