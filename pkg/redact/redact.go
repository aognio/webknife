// Package redact provides centralized secret detection and masking utilities.
//
// This package has no dependencies on other Webknife feature packages.
package redact

import "strings"

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

// IsSensitiveHeader returns true if the header name should be redacted.
func IsSensitiveHeader(name string) bool {
	return sensitiveHeaders[strings.ToLower(name)]
}

// IsSensitiveQueryParam returns true if the query param should be redacted.
func IsSensitiveQueryParam(name string) bool {
	return sensitiveQueryParams[strings.ToLower(name)]
}

// Value masks a sensitive string, keeping first 2 and last 2 characters.
func Value(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}
