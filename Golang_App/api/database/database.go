package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	MongoClient *mongo.Client
}

func ConnectDatabase() Db {

	log.Println("Database connecting...")
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_DB_URL"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return Db{}
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected.")
	return Db{MongoClient: client}

}
