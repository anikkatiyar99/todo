package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClient is a function that initializes a connection with the database
func GetClient() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL")) // Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Not Connected to MongoDB!")
	}
	return client
}

//Client Database instance
var Client *mongo.Client = GetClient()

//OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("Cluster0").Collection(collectionName)

	return collection
}
