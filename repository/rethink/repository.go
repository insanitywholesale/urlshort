package rethink

import (
	"urlshort/shortener"
)

type rethinkRepo struct {
}

func newRethinkClient() {
}

func NewRethinkRepo() {
}

func (r *rethinkRepo) Find(code string) (*shortener.Redirect, error) {
}

func (r *rethinkRepo) Store(redirect *shortener.Redirect) error {
}
