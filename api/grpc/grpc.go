package grpcapi

import (
	"urlshort/shortener"
)

type RedirectHandler interface {
	PostRedir(string) (*shortener.Redirect, error)
}

type handler struct {
	redirectService shortener.RedirectService
}

type Redirect struct{}

func NewHandlerGRPC(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

func (r *Redirect) grpcToRedirect(link string) *shortener.Redirect {
	redirect := &shortener.Redirect{}
	redirect.URL = link
	return redirect
}

func (h *handler) PostRedir(ll string) (*shortener.Redirect, error) {
	redir := &Redirect{}
	r := redir.grpcToRedirect(ll)
	h.redirectService.Store(r)
	return r, nil
}
