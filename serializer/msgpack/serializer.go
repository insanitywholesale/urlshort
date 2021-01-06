package msgpack

import (
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
	"log"
	"urlshort/shortener"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	log.Print("input:", input)
	err := msgpack.Unmarshal(input, redirect)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	log.Print("rawMsg:", rawMsg)
	return rawMsg, nil
}
