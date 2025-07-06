package connection

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongodbConn() (*mongo.Client, context.Context, context.CancelFunc) {
	mongodbUrl := "mongodb://" + os.Getenv("MONGODB_HOST")
	clientOptions := options.Client().ApplyURI(mongodbUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx, cancel
}
