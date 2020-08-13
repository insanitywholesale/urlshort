package mock

import (
	"fmt"
	"urlshort/shortener"
)

// Type mockRepo is a list of storedRedirect
type mockRepo []*shortener.Redirect

// List of storedRedirect is of data type mockRepo
var redirectList mockRepo = []*shortener.Redirect{}

func mkRedirList() mockRepo {
	//var redirectList[0] storedRedirects
	var redirList mockRepo = []*shortener.Redirect{
		&shortener.Redirect{
			Code:      "0123",
			URL:       "http://git.example.com",
			CreatedAt: 3254,
		},
		&shortener.Redirect{
			Code:      "1234",
			URL:       "https://distro.watch",
			CreatedAt: 7873,
		},
	}
	return redirList
}

func NewMockRepo() (shortener.RedirectRepo, error) {
	redirectList = mkRedirList()
	repo := &mockRepo{}
	return repo, nil
}

func (r *mockRepo) Find(code string) (*shortener.Redirect, error) {
	for _, redir := range redirectList {
		if redir.Code == code {
			return redir, nil
		}
	}
	return nil, fmt.Errorf("Code not found")
}

func (r *mockRepo) Store(sr *shortener.Redirect) error {
	redirectList = append(redirectList, sr)
	return nil
}
