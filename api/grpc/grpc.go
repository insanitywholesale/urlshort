package grpcapi

import (
	"context"
	"log"
	protos "gitlab.com/insanitywholesale/urlshort/proto/shorten"
	"gitlab.com/insanitywholesale/urlshort/shortener"
)

type ShortenRequest struct {
	link string
}

var redirSrv shortener.RedirectService

func NewHandlerGRPC(redirectService shortener.RedirectService) {
	//provide access to redirectService create in main.go
	redirSrv = redirectService
}

//figure out how to have this func attached both to `ShortenRequest` and `handler`
func (sr *ShortenRequest) GetShortURL(ctx context.Context, ll *protos.LongLink) (*protos.ShortLink, error) {
	log.Println("in GetShortURL")
	r := &shortener.Redirect{URL: ll.Link}
	err := redirSrv.Store(r)
	if err != nil {
		log.Println("oofie", err)
		return nil, err
	}
	return &protos.ShortLink{Link: r.Code}, nil
}
