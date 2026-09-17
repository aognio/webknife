// Package auth provides HTTP Basic Authentication middleware.
//
// This package has no dependencies on other Webknife feature packages.
package auth

import (
	"crypto/subtle"
	"net/http"
)

// Config holds configuration for Basic Authentication.
type Config struct {
	Username string
	Password string
	Realm    string
}

// Basic returns middleware that enforces HTTP Basic Authentication.
func Basic(cfg Config, next http.Handler) http.Handler {
	if cfg.Realm == "" {
		cfg.Realm = "Restricted"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+cfg.Realm+`"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		usernameMatch := subtle.ConstantTimeCompare([]byte(user), []byte(cfg.Username)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(cfg.Password)) == 1
		if !usernameMatch || !passwordMatch {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+cfg.Realm+`"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
