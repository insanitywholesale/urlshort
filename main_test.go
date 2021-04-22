package main

import (
	"bytes"
	"context"
	"github.com/vmihailenco/msgpack/v5"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	hg "urlshort/api/grpc"
	shortie "urlshort/api/grpc"
	protos "urlshort/proto/shorten"
	"urlshort/shortener"
)

// Test is still incomplete
// Should be integration test maybe
func TestGet(t *testing.T) {
	// initialize mock repo
	repo, err := chooseRepo()
	if err != nil {
		t.Error("repo oopsie")
	}
	// make redirect service
	service := shortener.NewRedirectService(repo)
	// create router based on the above service
	r := makeRouter(service)
	// create and start a test server
	testServer := httptest.NewServer(r)
	// do a simple Get request on preexisting redirect
	// said redirect can be found in the mock repo source
	res, err := http.Get(testServer.URL + "/1234")
	// be responsible and close the response body
	res.Body.Close()
	if err != nil {
		t.Error("tfw GET error:", err)
	}
	// close the test server
	testServer.Close()
}

func TestPostJSON(t *testing.T) {
	// initialize mock repo
	repo, err := chooseRepo()
	if err != nil {
		t.Error("repo oopsie")
	}
	// make redirect service
	service := shortener.NewRedirectService(repo)
	// create router based on the above service
	r := makeRouter(service)
	// create and start a test server
	testServer := httptest.NewServer(r)

	// create some data in the form of an io.Reader from a string of json
	jsonData := []byte(`{"url": "https://todo.distro.watch"}`)
	// do a simple Post request with the above data
	res, err := http.Post(testServer.URL, "application/json", bytes.NewBuffer(jsonData))
	// be responsible and close the response body
	res.Body.Close()
	if err != nil {
		t.Error("tfw POST error:", err)
	}
	// close the test server
	testServer.Close()
}

func TestPostXML(t *testing.T) {
	// initialize mock repo
	repo, err := chooseRepo()
	if err != nil {
		t.Error("repo oopsie")
	}
	// make redirect service
	service := shortener.NewRedirectService(repo)
	// create router based on the above service
	r := makeRouter(service)
	// create and start a test server
	testServer := httptest.NewServer(r)

	// create some data in the form of an io.Reader from a string of xml
	xmlData := []byte(`<redirect><url>https://distro.watch</url></redirect>`)
	// do a simple Post request with the above data
	res, err := http.Post(testServer.URL, "application/xml", bytes.NewBuffer(xmlData))
	// be responsible and close the response body
	res.Body.Close()
	if err != nil {
		t.Error("tfw POST error:", err)
	}
	// close the test server
	testServer.Close()
}

func TestPostMsgPack(t *testing.T) {
	// initialize mock repo
	repo, err := chooseRepo()
	if err != nil {
		t.Error("repo oopsie")
	}
	// make redirect service
	service := shortener.NewRedirectService(repo)
	// create router based on the above service
	r := makeRouter(service)
	// create and start a test server
	testServer := httptest.NewServer(r)

	// create some data in the form of an io.Reader from a string of xml
	msgpackData, err := msgpack.Marshal(&shortener.Redirect{URL: "https://inherently.xyz"})
	// do a simple Post request with the above data
	res, err := http.Post(testServer.URL, "application/x-msgpack", bytes.NewBuffer(msgpackData))
	// be responsible and close the response body
	res.Body.Close()
	if err != nil {
		t.Error("tfw POST error:", err)
	}
	// close the test server
	testServer.Close()
}

func TestGPRC(t *testing.T) {
	const bufSize = 1024 * 1024
	var listener *bufconn.Listener
	listener = bufconn.Listen(bufSize)
	grpcServer := ggrpc.NewServer()
	protos.RegisterShortenRequestServer(grpcServer, &shortie.ShortenRequest{})
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// initialize mock repo
	repo, err := chooseRepo()
	if err != nil {
		t.Error("repo oopsie")
	}
	// make redirect service
	service := shortener.NewRedirectService(repo)
	//http service needs to be up in order to test if the shortcode works
	// create router based on the above service
	r := makeRouter(service)
	// create and start a test server
	testServer := httptest.NewServer(r)

	// set up grpc connection
	ctx := context.Background()
	conn, err := ggrpc.DialContext(
		ctx,
		"bufnet",
		ggrpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		ggrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	// make handler for grpc
	handlergrpc := hg.NewHandlerGRPC(service)
	shortie.SaveHandler(handlergrpc)

	// make a grpc client
	client := protos.NewShortenRequestClient(conn)
	response, err := client.GetShortURL(ctx, &protos.LongLink{Link: "https://distro.watch"})
	// test if shortcode works
	http.Get(testServer.URL + response.Link)

	// close test http server
	testServer.Close()
}
