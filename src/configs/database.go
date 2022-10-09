package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// Database Config
	clientOptions := options.Client().ApplyURI(EnvMongoURI())
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Cancel context to avoid memory leak
	defer cancel()

	// Ping our db connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// Getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	collection := DB.Database("stanza").Collection(collectionName)
	return collection
}
