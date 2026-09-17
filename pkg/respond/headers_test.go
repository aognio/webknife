package respond_test

import (
	"testing"

	"github.com/webknife/webknife/pkg/respond"
)

func TestParseHeaders(t *testing.T) {
	headers, err := respond.ParseHeaders([]string{"X-Foo: bar", "X-Baz: qux"})
	if err != nil {
		t.Fatal(err)
	}
	if headers["X-Foo"] != "bar" {
		t.Fatalf("expected X-Foo=bar, got %s", headers["X-Foo"])
	}
	if headers["X-Baz"] != "qux" {
		t.Fatalf("expected X-Baz=qux, got %s", headers["X-Baz"])
	}
}

func TestParseHeaders_Invalid(t *testing.T) {
	_, err := respond.ParseHeaders([]string{"invalid"})
	if err == nil {
		t.Fatal("expected error for invalid header")
	}
	_, err = respond.ParseHeaders([]string{": value"})
	if err == nil {
		t.Fatal("expected error for empty header name")
	}
}
