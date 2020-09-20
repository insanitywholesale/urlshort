package main

import (
	"net/http"
	"testing"
)

// Test is still inclomplete
// Find way to start and stop service
// Should be integration test
func TestGet(t *testing.T) {
	val, err := http.Get("http://localhost:8000/1234")
	if err != nil {
		t.Errorf("no respond")
	}
	t.Log("val:", val)
}
