package webknife

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

const DefaultMaxBodySize = int64(1 << 20) // 1MB

type EchoConfig struct {
	MaxBodySize int64
}

type RequestInspection struct {
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

func NewEchoHandler(cfg EchoConfig) http.Handler {
	if cfg.MaxBodySize <= 0 {
		cfg.MaxBodySize = DefaultMaxBodySize
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inspection := inspectRequest(r, cfg.MaxBodySize)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(inspection)
	})
}

func inspectRequest(r *http.Request, maxBodySize int64) *RequestInspection {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	inspection := &RequestInspection{
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
		inspection.Headers[k] = v
	}

	if r.TLS != nil {
		inspection.TLS = &TLSInfo{
			Version:    tlsVersionString(r.TLS.Version),
			ServerName: r.TLS.ServerName,
		}
	}

	for _, cookie := range r.Cookies() {
		inspection.Cookies = append(inspection.Cookies, &CookieInfo{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	if r.Body != nil {
		bodyBytes, truncated, isBinary := readBody(r.Body, maxBodySize)
		if isBinary {
			inspection.BodyBinary = true
			inspection.Body = encodeBinaryBody(bodyBytes)
		} else {
			inspection.Body = string(bodyBytes)
		}
		inspection.BodyTruncated = truncated
	}

	return inspection
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
					strings.Replace(
						string([]byte{b}),
						`\`, `\\`, -1,
					),
					"\\x",
				),
				"\\x", "", -1,
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
