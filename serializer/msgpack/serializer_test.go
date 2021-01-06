package msgpack

import (
	"testing"
	"urlshort/shortener"
)

var localRedir Redirect

func TestJson(t *testing.T) {
	// sample redirect
	shortRedir := &shortener.Redirect{
		Code:      "roastBread",
		URL:       "http://roast.bread",
		CreatedAt: 20202020,
	}
	// encode shortener redirect struct to msgpack
	encodedRedirect, err := localRedir.Encode(shortRedir)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	t.Log("encodedRedirect:", encodedRedirect)

	// no idea how to put a string that doesn't result in
	// msgpack: invalid code=5b decoding map length
	// so just let it decode what it encoded
	msgpackBytes := encodedRedirect
	// decode msgpack into shortener redirect struct
	decodedRedirect, err := localRedir.Decode(msgpackBytes)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	t.Log("decodedRedirect:", decodedRedirect)
}
