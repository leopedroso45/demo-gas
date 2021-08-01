package mongo

import (
	"context"
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

func CreateUser(newUser model.User) (interface{}, error) {
	client, ctx, err := createClient()
	if err != nil {
		log.Panicln("Error creating mongodb client", err)
		return nil, err
	}
	demogasDb := client.Database("demogas")
	usersCollection := demogasDb.Collection("users")

	emailCheck := usersCollection.FindOne(ctx, bson.M{"email": newUser.Email})
	usernameCheck := usersCollection.FindOne(ctx, bson.M{"username": newUser.Username})
	if emailCheck.Err() == nil {
		return nil, fmt.Errorf("user with email: %s already exists", newUser.Email)
	} else if usernameCheck.Err() == nil {
		return nil, fmt.Errorf("user with username: %s already exists", newUser.Username)
	}
	newUser.PassToSha1()
	result, err := usersCollection.InsertOne(ctx, newUser.ToBSON())
	if err != nil {
		return nil, fmt.Errorf("error inserting a new document - %s", err)
	}
	fmt.Println("Result:", result.InsertedID)
	defer client.Disconnect(ctx)
	return result.InsertedID, nil
}

func DeleteUser(newUser model.User) ([]byte, error) {
	client, ctx, err := createClient()
	if err != nil {
		log.Panicln("Error creating mongodb client", err)
		return nil, err
	}
	demogasDb := client.Database("demogas")
	usersCollection := demogasDb.Collection("users")

	newUser.PassToSha1()
	result := usersCollection.FindOneAndDelete(ctx, bson.M{"email": newUser.Email, "password": newUser.Password})

	deletedUser, err := result.DecodeBytes()

	if err != nil {
		fmt.Println("Something went wrong decoding the deleted user to bson ", err)
		return deletedUser, err
	}

	defer client.Disconnect(ctx)
	return deletedUser, nil
}

func EditUser(newUser model.User) ([]byte, error) {
	client, ctx, err := createClient()
	if err != nil {
		log.Panicln("Error creating mongodb client", err)
		return nil, err
	}
	demogasDb := client.Database("demogas")
	usersCollection := demogasDb.Collection("users")

	newUser.PassToSha1()
	result := usersCollection.FindOneAndUpdate(ctx, bson.M{"email": newUser.Email}, bson.M{"$set": bson.M{"name": newUser.Name, "password": newUser.Password, "username": newUser.Username}})

	if result.Err() != nil {
		fmt.Println("Something went wrong updating data ", err)
		return nil, err
	}

	resultBson, _ := result.DecodeBytes()

	defer client.Disconnect(ctx)
	return resultBson, nil
}
