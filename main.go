package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/leonardom/go-url-shortener/api"
	mr "github.com/leonardom/go-url-shortener/repository/mongo"
	rr "github.com/leonardom/go-url-shortener/repository/redis"
	"github.com/leonardom/go-url-shortener/shortener"
)

// https://www.google.com -> 98sj1-293
// https://localhost.com:8000/98sj1-293 -> https://www.google.com

// repos <- service -> serializer -> http

func main()  {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func ()  {
		port := httpPort()
		fmt.Printf("Listening on port %s\n", port)
		errs <- http.ListenAndServe(port, router)
	}()

	go func ()  {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <- c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("POST") != "" {
		port = os.Getenv("POST")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("DB_TYPE") {
	case "redis":
		redisURL :=  os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}