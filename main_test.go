package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshort/shortener"
)

// Test is still incomplete
// Should be integration test maybe
func TestGet(t *testing.T) {
	// initialize mock repo
	repo, repoErr := chooseRepo()
	if repoErr != nil {
		t.Errorf("repo oopsie")
	}
	// make redirect service
	service := shortener.NewRedirectService(repo)
	// create router based on the above service
	r := makeRouter(service)
	// create and start a test server
	testServer := httptest.NewServer(r)

	// do a simple Get request on preexisting redirect
	// said redirect can be found in the mock repo source
	res, err := http.Get(testServer.URL + "/1234")
	// be responsible and close the response body
	res.Body.Close()
	if err != nil {
		t.Errorf("tfw no respond")
	}
	// close the test server
	testServer.Close()
}
