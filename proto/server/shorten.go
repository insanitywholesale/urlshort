package server

import (
	"context"
	"log"
	h "urlshort/api/grpc"
	protos "urlshort/proto/shorten"
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
