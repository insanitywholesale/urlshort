package grpcapi

import (
	"context"
	"log"
	protos "urlshort/proto/shorten"
	"urlshort/shortener"
)

type ShortenRequest struct {
	link string
}

var redirSrv shortener.RedirectService

func (sr *ShortenRequest) GetShortURL(ctx context.Context, ll *protos.LongLink) (*protos.ShortLink, error) {
	r := &shortener.Redirect{URL: ll.Link}
	err := redirSrv.Store(r)
	if err != nil {
		log.Println("oofie", err)
		return &protos.ShortLink{}, err
	}
	return &protos.ShortLink{Link: r.Code}, nil
}

func NewHandlerGRPC(redirectService shortener.RedirectService) {
	redirSrv = redirectService
}
