package webknife

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type RespondConfig struct {
	StatusCode  int
	Body        string
	BodyFile    string
	ContentType string
	Headers     map[string]string
}

func NewRespondHandler(cfg RespondConfig) (http.Handler, error) {
	if cfg.StatusCode < 100 || cfg.StatusCode > 599 {
		return nil, fmt.Errorf("invalid status code: %d", cfg.StatusCode)
	}

	var body []byte
	if cfg.BodyFile != "" {
		data, err := os.ReadFile(cfg.BodyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read body file: %w", err)
		}
		body = data
	} else {
		body = []byte(cfg.Body)
	}

	contentType := cfg.ContentType
	if contentType == "" && len(body) > 0 {
		contentType = http.DetectContentType(body)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range cfg.Headers {
			w.Header().Set(k, v)
		}
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(cfg.StatusCode)
		w.Write(body)
	}), nil
}

func ParseHeaders(headerFlags []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, h := range headerFlags {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format: %s (expected 'Name: Value')", h)
		}
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if name == "" {
			return nil, fmt.Errorf("empty header name in: %s", h)
		}
		headers[name] = value
	}
	return headers, nil
}
