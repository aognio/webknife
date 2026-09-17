package webknife

import (
	"net/http"
	"path/filepath"
	"strings"
)

type StaticConfig struct {
	Root string
}

func NewStaticHandler(cfg StaticConfig) (http.Handler, error) {
	absRoot, err := filepath.Abs(cfg.Root)
	if err != nil {
		return nil, err
	}
	return http.StripPrefix("/", http.FileServer(http.Dir(absRoot))), nil
}

func StaticServeMux(root string) (*http.ServeMux, error) {
	handler, err := NewStaticHandler(StaticConfig{Root: root})
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	return mux, nil
}

func SanitizePath(path string) string {
	path = filepath.Clean(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
