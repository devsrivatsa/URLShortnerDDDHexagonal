package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	api "github.com/devsrivatsa/URLShortnerDDDHexagonal/api"
	mongoRepo "github.com/devsrivatsa/URLShortnerDDDHexagonal/repository/mongodb"
	redisRepo "github.com/devsrivatsa/URLShortnerDDDHexagonal/repository/redis"
	"github.com/devsrivatsa/URLShortnerDDDHexagonal/urlShortner"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func chooseRepo() urlShortner.RedirectRepository {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	repoType := os.Getenv("REPO_TYPE")
	if repoType == "" {
		log.Fatal("REPO_TYPE is not set")
	}
	switch repoType {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redisRepo.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal("Error creating redis repository")
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongoRepo.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal("Error creating mongo repository")
		}
		return repo
	default:
		log.Fatal("Invalid REPO_TYPE")
	}
	return nil
}

func main() {
	repo := chooseRepo()
	service := urlShortner.New(repo)
	handler := api.NewRedirectHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errsChannel := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errsChannel <- http.ListenAndServe(":8000", r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errsChannel <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated , errs: %s\n", <-errsChannel)
}
