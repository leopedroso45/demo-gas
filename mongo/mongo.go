package mongo

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"

	"demogas.com/m/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createClient() (*mongo.Client, context.Context, error) {
	var cred options.Credential

	cred.AuthSource = "admin"
	cred.Username = "root"
	cred.Password = "password"
	//mongoURI := os.Getenv("MONGO_URI")
	mongoURI := "mongodb://localhost:27017"
	fmt.Println("--------mongo" + mongoURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI).SetAuth(cred))
	if err != nil {
		log.Printf(`MongoDB can't create the client %s `, err)
		return client, nil, err
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Printf(`MongoDB can't connect %s `, err)
		return client, nil, err
	}
	return client, ctx, nil
}

func CreateUser(newUser model.User) error {
	client, ctx, err := createClient()
	if err != nil {
		fmt.Println("--------mongo111111111")
		log.Panicln(err)
		return err
	}
	demogasDb := client.Database("demogas")
	usersCollection := demogasDb.Collection("users")

	result, err := usersCollection.InsertOne(ctx, bson.M{"_id": newUser.Id, "name": newUser.Name, "username": newUser.Username, "email": newUser.Email, "password": toSha1(newUser.Password)})
	if err != nil {
		fmt.Println("--------mongo2222222")
		log.Panicln(err)
		return err
	}
	fmt.Println("Result:", result)
	defer client.Disconnect(ctx)
	return nil
}

func toSha1(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	return string(bs)
}
