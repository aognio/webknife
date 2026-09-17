package webknife

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyConfig struct {
	Upstream string
}

func NewProxyHandler(cfg ProxyConfig) (http.Handler, error) {
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
