package redis

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/devsrivatsa/URLShortnerDDDHexagonal/urlShortner"
	"github.com/go-redis/redis"
	errs "github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	// Retry logic for Redis connection
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		_, err = client.Ping().Result()
		if err == nil {
			return client, nil
		}
		log.Printf("Redis connection failed, retrying... (%d/%d)", i+1, maxRetries)
		time.Sleep(time.Duration(i+2) * time.Second)
	}

	return client, err
}

func NewRedisRepository(redisURL string) (*redisRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client

	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*urlShortner.Redirect, error) {
	key := r.generateKey(code)
	log.Printf("Attempting to find key: %s", key)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		log.Printf("Redis error during find: %v", err)
		return nil, errs.Wrap(err, "repository.redisRepository.Find")
	}
	log.Printf("Redis returned data: %+v", data)
	if len(data) == 0 {
		log.Printf("No data found for key: %s", key)
		return nil, urlShortner.ErrRedirectNotFound
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errs.Wrap(err, "repository.redisRepository.Find")
	}
	redirect := &urlShortner.Redirect{
		Code:      data["code"],
		URL:       data["url"],
		CreatedAt: createdAt,
	}
	return redirect, nil
}

func (r *redisRepository) Store(redirect *urlShortner.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"url":        redirect.URL,
		"code":       redirect.Code,
		"created_at": redirect.CreatedAt,
	}
	log.Printf("Attempting to store key: %s, data: %+v", key, data)
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		log.Printf("Error storing to Redis: %v", err)
		return errs.Wrap(err, "repository.redisRepository.Store")
	}
	log.Printf("Successfully stored redirect with code: %s", redirect.Code)
	return nil
}
