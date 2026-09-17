package webknife

import (
	"strings"
)

var sensitiveHeaders = map[string]bool{
	"authorization":       true,
	"cookie":              true,
	"set-cookie":          true,
	"proxy-authorization": true,
}

var sensitiveQueryParams = map[string]bool{
	"token":    true,
	"key":      true,
	"secret":   true,
	"password": true,
	"api_key":  true,
	"apikey":   true,
}

func RedactHeader(name string) bool {
	return sensitiveHeaders[strings.ToLower(name)]
}

func RedactQueryParam(name string) bool {
	return sensitiveQueryParams[strings.ToLower(name)]
}

func RedactValue(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}
