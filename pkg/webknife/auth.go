package webknife

import (
	"crypto/subtle"
	"net/http"
)

type BasicAuthConfig struct {
	Username string
	Password string
	Realm    string
}

func BasicAuth(cfg BasicAuthConfig, next http.Handler) http.Handler {
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
