package server

import (
	"context"
	"log"
	protos "urlshort/proto/shorten"
	h "urlshort/api/grpc"
)

type ShortenRequest struct {
	link string
}

var handler h.RedirectHandler

func SaveHandler(h h.RedirectHandler) {
	handler = h
}

func (sr *ShortenRequest) GetShortURL(ctx context.Context, ll *protos.LongLink) (*protos.ShortLink, error) {
	r, err := handler.PostRedir(ll.Link)
	if err != nil {
		log.Fatal(err)
	}

	return &protos.ShortLink{Link: r.Code}, nil
}
