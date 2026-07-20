package utils

import "testing"

func TestHTTPInvalidRequestReturnsError(t *testing.T) {
	_, err := NewHttp().Get("://invalid-url")
	if err == nil {
		t.Fatal("Get() error = nil, want invalid URL error")
	}
}
