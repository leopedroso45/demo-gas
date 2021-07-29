package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateClient() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf(`MongoDB can't create the client %s `, err)
		return client, err
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Printf(`MongoDB can't connect %s `, err)
		return client, err
	}
	defer client.Disconnect(ctx)
	return client, nil
}
