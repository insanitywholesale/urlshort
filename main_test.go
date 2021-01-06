package main

import (
	"net/http"
	"urlshort/shortener"
	"testing"
	"time"
)

// Test is still incomplete
// Should be integration test
// This sometimes works as it should
// There is definitely a race condition here
func TestGet(t *testing.T) {
	repo, repoErr := chooseRepo()
	if repoErr != nil {
		t.Errorf("repo oopsie")
	}
	service := shortener.NewRedirectService(repo)
	go setupHTTP(service)
	time.Sleep(5 * time.Second)
	val, err := http.Get("http://localhost:8000/1234")
	if err != nil {
		t.Errorf("no respond")
	}
	t.Log("val:", val)
}
