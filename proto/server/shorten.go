package server

import (
	"context"
	protos "urlshort/proto/shorten"
)

type ShortenRequest struct {
	link string
}

func (sr *ShortenRequest) GetShortURL(ctx context.Context, ll *protos.LongLink) (*protos.ShortLink, error) {
	return &protos.ShortLink{Link: "got_it"}, nil
}
