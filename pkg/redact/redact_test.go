package redact_test

import (
	"testing"

	"github.com/aognio/webknife/pkg/redact"
)

func TestRedactValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"a", "****"},
		{"ab", "****"},
		{"abc", "****"},
		{"abcd", "****"},
		{"abcde", "ab*de"},
		{"abcdef", "ab**ef"},
		{"secret123", "se*****23"},
	}
	for _, tt := range tests {
		got := redact.Value(tt.input)
		if got != tt.want {
			t.Errorf("RedactValue(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRedactHeader(t *testing.T) {
	if !redact.IsSensitiveHeader("Authorization") {
		t.Fatal("expected Authorization to be sensitive")
	}
	if !redact.IsSensitiveHeader("authorization") {
		t.Fatal("expected authorization to be sensitive")
	}
	if !redact.IsSensitiveHeader("Cookie") {
		t.Fatal("expected Cookie to be sensitive")
	}
	if redact.IsSensitiveHeader("X-Custom") {
		t.Fatal("expected X-Custom to not be sensitive")
	}
}
