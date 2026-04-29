package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongoDB creates and pings a MongoDB client with a timeout.
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// Limit connection/ping operations to a short timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		// Log fatal to surface config/connect issues early.
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		// Fail fast if the server is unreachable.
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")
	return client, nil
}
