package mock

import (
	"gitlab.com/insanitywholesale/urlshort/shortener"
	"testing"
)

func TestMock(t *testing.T) {
	// initialize repo
	repo, err := NewMockRepo()
	if err != nil {
		t.Log("error creating repo")
	}
	// try finding one of the predefined redirects
	redirCode := "1234"
	findRedir, err := repo.Find(redirCode)
	// check if the Find method worked
	if err != nil {
		t.Log("error finding redirect:", err)
	}
	// check if the redirect returned is correct
	if findRedir.URL != "https://distro.watch" || findRedir.CreatedAt != 7873 {
		t.Log("error with the returned redirect")
	}
	// make a redirect to store
	storeRedir := &shortener.Redirect{
		Code:      "_o+9z-2",
		URL:       "https://jvt.me",
		CreatedAt: 56382910,
	}
	// try storing the redirect
	err = repo.Store(storeRedir)
	// check if the Store method worked
	if err != nil {
		t.Log("error storing redirect:", err)
	}
}
