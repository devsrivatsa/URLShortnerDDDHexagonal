package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/devsrivatsa/URLShortnerDDDHexagonal/urlShortner"
	errs "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Println("MongoDB connection failed..")
		return nil, errs.Wrap(err, "repository.NewMongoClient")
	}

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = client.Ping(ctx, readpref.Primary())
		if err == nil {
			log.Printf("MongoDB connection established..")
			return client, nil
		}
		log.Printf("MongoDB connection failed, retrying... (%d/%d)", i+1, maxRetries)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return client, err

}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (*MongoRepository, error) {
	repo := &MongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewMongoRepository")
	}
	repo.client = client

	return repo, nil
}

func (r *MongoRepository) Find(code string) (*urlShortner.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	redirect := &urlShortner.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.Wrap(urlShortner.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errs.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil

}

func (r *MongoRepository) Store(redirect *urlShortner.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
