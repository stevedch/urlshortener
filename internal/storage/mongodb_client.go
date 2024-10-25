package storage

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// MongoDBClient is the interface that defines MongoDB connection methods
type MongoDBClient interface {
	Connect(uri string) rxgo.Observable
	Disconnect() rxgo.Observable
	GetClient() *mongo.Client
}

// MongoDBService implements the MongoDBClient interface
type MongoDBService struct {
	Client *mongo.Client
}

// Connect establishes a reactive connection to MongoDB
func (m *MongoDBService) Connect(uri string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			ch <- rxgo.Error(err)
			return
		}

		if err = client.Ping(ctx, nil); err != nil {
			ch <- rxgo.Error(err)
			return
		}

		m.Client = client
		ch <- rxgo.Of(client)
		log.Println("Connected to MongoDB reactively")
	}})
}

// Disconnect closes the MongoDB connection reactively
func (m *MongoDBService) Disconnect() rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := m.Client.Disconnect(ctx); err != nil {
			ch <- rxgo.Error(err)
			return
		}

		ch <- rxgo.Of("Disconnected from MongoDB")
		log.Println("Disconnected from MongoDB reactively")
	}})
}

// GetClient provides access to the underlying MongoDB client
func (m *MongoDBService) GetClient() *mongo.Client {
	return m.Client
}
