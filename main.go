package main

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	shortie "urlshort/api/grpc"
	hh "urlshort/api/http"
	protos "urlshort/proto/shorten"
	mockrepo "urlshort/repository/mock"
	mr "urlshort/repository/mongo"
	rr "urlshort/repository/redis"
	"urlshort/shortener"
)

func grpcPort() string {
	port := "4040"
	if os.Getenv("GRPC_PORT") != "" {
		port = os.Getenv("GRPC_PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() (shortener.RedirectRepo, error) {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepo(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo, nil
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepo(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo, nil
	default:
		repo, err := mockrepo.NewMockRepo()
		if err != nil {
			log.Fatal(err)
		}
		return repo, nil
	}
	return nil, nil
}

func setupGRPC(servicegrpc shortener.RedirectService) *grpc.Server {
	grpcServer := grpc.NewServer()
	protos.RegisterShortenRequestServer(grpcServer, &shortie.ShortenRequest{})
	shortie.NewHandlerGRPC(servicegrpc)
	return grpcServer
}

func makeRouter(service shortener.RedirectService) http.Handler {
	handler := hh.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)
	return r
}

func setupHTTP(service shortener.RedirectService) http.Handler {
	r := makeRouter(service)
	return r
}

func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	repo, repoErr := chooseRepo()
	if repoErr != nil {
		fmt.Println("No database backend has been selected")
		os.Exit(1)
	}

	service := shortener.NewRedirectService(repo)

	/*
		go setupGRPC(service)
		defer setupHTTP(service)
	*/
	//gS := setupGRPC(service)
	hH := setupHTTP(service)
	gS := setupGRPC(service)
	log.Fatal(http.ListenAndServe(":8086", httpGrpcRouter(gS, hH)))
}
