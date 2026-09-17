// Package echo provides a request inspection handler that reflects
// incoming HTTP requests as structured JSON.
//
// This package has no dependencies on other Webknife feature packages.
package echo

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

// DefaultMaxBodySize is the default maximum body size for inspection (1MB).
const DefaultMaxBodySize = int64(1 << 20)

// Config holds configuration for the echo handler.
type Config struct {
	// MaxBodySize is the maximum number of body bytes to retain.
	// Bodies exceeding this are truncated. Zero uses DefaultMaxBodySize.
	MaxBodySize int64
}

// Inspection represents a structured diagnostic view of an HTTP request.
type Inspection struct {
	Method        string              `json:"method"`
	Scheme        string              `json:"scheme"`
	Host          string              `json:"host"`
	Path          string              `json:"path"`
	Query         map[string][]string `json:"query"`
	Proto         string              `json:"proto"`
	ContentLength int64               `json:"content_length"`
	RemoteAddr    string              `json:"remote_addr"`
	Headers       map[string][]string `json:"headers"`
	URI           string              `json:"request_uri"`
	TLS           *TLSInfo            `json:"tls,omitempty"`
	Body          string              `json:"body"`
	BodyTruncated bool                `json:"body_truncated,omitempty"`
	BodyBinary    bool                `json:"body_binary,omitempty"`
	Cookies       []*CookieInfo       `json:"cookies,omitempty"`
}

type TLSInfo struct {
	Version    string `json:"version"`
	ServerName string `json:"server_name"`
}

type CookieInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// New returns an http.Handler that reflects requests as JSON.
func New(cfg Config) http.Handler {
	maxBody := cfg.MaxBodySize
	if maxBody <= 0 {
		maxBody = DefaultMaxBodySize
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inspection := inspectRequest(r, maxBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(inspection)
	})
}

func inspectRequest(r *http.Request, maxBodySize int64) *Inspection {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	ins := &Inspection{
		Method:        r.Method,
		Scheme:        scheme,
		Host:          r.Host,
		Path:          r.URL.Path,
		Query:         r.URL.Query(),
		Proto:         r.Proto,
		ContentLength: r.ContentLength,
		RemoteAddr:    r.RemoteAddr,
		URI:           r.RequestURI,
		Headers:       make(map[string][]string),
	}

	for k, v := range r.Header {
		ins.Headers[k] = v
	}

	if r.TLS != nil {
		ins.TLS = &TLSInfo{
			Version:    tlsVersionString(r.TLS.Version),
			ServerName: r.TLS.ServerName,
		}
	}

	for _, cookie := range r.Cookies() {
		ins.Cookies = append(ins.Cookies, &CookieInfo{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	if r.Body != nil {
		bodyBytes, truncated, isBinary := readBody(r.Body, maxBodySize)
		if isBinary {
			ins.BodyBinary = true
			ins.Body = encodeBinaryBody(bodyBytes)
		} else {
			ins.Body = string(bodyBytes)
		}
		ins.BodyTruncated = truncated
	}

	return ins
}

func readBody(body io.Reader, maxSize int64) (data []byte, truncated bool, isBinary bool) {
	limitedReader := io.LimitReader(body, maxSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, false, false
	}
	if int64(len(data)) > maxSize {
		data = data[:maxSize]
		truncated = true
	}
	isBinary = !utf8.Valid(data)
	return data, truncated, isBinary
}

func encodeBinaryBody(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		if b >= 32 && b < 127 && b != '\\' && b != '"' {
			sb.WriteByte(b)
		} else {
			sb.WriteString(`\x`)
			sb.WriteString(strings.ToUpper(strings.Replace(
				strings.TrimPrefix(
					strings.Replace(string([]byte{b}), `\`, `\\`, -1),
					`\x`,
				),
				`\x`, "", -1,
			)))
		}
	}
	return sb.String()
}

func tlsVersionString(version uint16) string {
	switch version {
	case 0x0301:
		return "TLS 1.0"
	case 0x0302:
		return "TLS 1.1"
	case 0x0303:
		return "TLS 1.2"
	case 0x0304:
		return "TLS 1.3"
	default:
		return "unknown"
	}
}
