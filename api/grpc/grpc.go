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
type RedirectHandler interface {
	//figure this out
}

//copy of what is in http.go api
type handler struct {
	redirectService shortener.RedirectService
}

//might not be required -- doesn't even do anything at this point
var redirSrv shortener.RedirectService

func NewHandlerGRPC(redirectService shortener.RedirectService) {
	//probably need to add more to this
	redirSrv = redirectService
}

//figure out how to have this func attached both to `ShortenRequest` and `handler`
func (sr *ShortenRequest) GetShortURL(ctx context.Context, ll *protos.LongLink) (*protos.ShortLink, error) {
	log.Println("in GetShortURL")
	r := &shortener.Redirect{URL: ll.Link}
	err := redirSrv.Store(r)
	if err != nil {
		log.Println("oofie", err)
		return &protos.ShortLink{}, err
	}
	return &protos.ShortLink{Link: r.Code}, nil
}

