package json

import (
	"testing"
	"urlshort/shortener"
)

var localRedir Redirect

func TestJson(t *testing.T) {
	jsonBytes := []byte(`{"code":"freeDrinks", "url":"http://example.com/drink", "created_at": 9223372036854775007}`)
	decodedRedirect, err := localRedir.Decode(jsonBytes)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	t.Log("decodedRedirect:", decodedRedirect)

	shortRedir := &shortener.Redirect{
		Code:      "bakedBread",
		URL:       "http://bakesome.bread",
		CreatedAt: 9111111111111111112,
	}
	encodedRedirect, err := localRedir.Encode(shortRedir)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	t.Log("encodedRedirect:", encodedRedirect)
}
