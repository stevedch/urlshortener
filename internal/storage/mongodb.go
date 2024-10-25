package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectMongoDB establishes the connection to the MongoDB server
func ConnectMongoDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	log.Println("Connected to MongoDB")
	client = c
	return client
}

// DisconnectMongoDB disconnects the MongoDB client
func DisconnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}
