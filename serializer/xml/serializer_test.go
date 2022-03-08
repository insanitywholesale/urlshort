package xml

import (
	"gitlab.com/insanitywholesale/urlshort/shortener"
	"testing"
)

var localRedir Redirect

func TestJson(t *testing.T) {
	xmlBytes := []byte(`<?xml version="1.0" encoding="UTF-8"?><redirect><url>https://distro.watch</url><code>_flaming+delta-Pr0ject</code><createdat>192168213</createdat></redirect>`)
	decodedRedirect, err := localRedir.Decode(xmlBytes)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	t.Log("decodedRedirect:", decodedRedirect)

	shortRedir := &shortener.Redirect{
		Code:      "breadSticks",
		URL:       "http://sticks.bread",
		CreatedAt: 6660000000000000666,
	}
	encodedRedirect, err := localRedir.Encode(shortRedir)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	t.Log("encodedRedirect:", encodedRedirect)
}
