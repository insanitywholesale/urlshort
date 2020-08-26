package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	h "urlshort/api"
	protos "urlshort/proto/shorten"
	shortie "urlshort/proto/server"
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

func chooseRepo() shortener.RedirectRepo {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepo(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepo(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	default:
		repo, err := mockrepo.NewMockRepo()
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}

func setupGRPC() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("ye dun goofed")
	}
	grpcServer := grpc.NewServer()
	protos.RegisterShortenRequestServer(grpcServer, &shortie.ShortenRequest{})
	grpcerrs := make(chan error, 2)
	go func() {
		fmt.Println("Listening for grpc on port :4040")
		grpcerrs <- grpcServer.Serve(listener)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		grpcerrs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s\n", <-grpcerrs)
}

func setupHTTP() {
	repo := chooseRepo()
	if repo == nil {
		fmt.Println("No database backend has been selected")
		os.Exit(1)
	}

	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	httperrs := make(chan error, 2)
	go func() {
		fmt.Println("Listening for http on port :8000")
		httperrs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		httperrs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s\n", <-httperrs)
}

func main() {

	go setupGRPC()
	defer setupHTTP()

}
