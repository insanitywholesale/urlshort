package server

import (
	"context"
	"log"
	"urlshort/shortener"
	mr "urlshort/repository/mock"
	protos "urlshort/proto/shorten"
)

type ShortenRequest struct {
	link string
}

type Redirect struct{}

var firstTime bool = true
var redirectService shortener.RedirectService
var repo shortener.RedirectRepo

func firstTimeSetup() {
	repo, err := mr.NewMockRepo()
	if err != nil {
		log.Fatal(err)
	}
	redirectService = shortener.NewRedirectService(repo)
}

func (r *Redirect) grpcToRedirect(link string) *shortener.Redirect {
	redirect := &shortener.Redirect{}
	redirect.URL = link
	return redirect
}

func (sr *ShortenRequest) GetShortURL(ctx context.Context, ll *protos.LongLink) (*protos.ShortLink, error) {
	if firstTime {
		firstTimeSetup()
		firstTime = false
	}
	redirStruct := &Redirect{}
	r := redirStruct.grpcToRedirect(ll.Link)
	err := redirectService.Store(r)
	if err != nil {
		log.Fatal(err)
	}

	return &protos.ShortLink{Link: r.Code}, nil
}
