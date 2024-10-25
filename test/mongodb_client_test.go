package test

import (
	"context"
	"testing"
	"time"
	"urlshortener/internal/storage"

	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestMongoDBServiceConnect tests the Connect method of MongoDBService
func TestMongoDBServiceConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock MongoDB client
	mockClient := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	// Instantiate MongoDBService with mock client
	service := &storage.MongoDBService{
		Client: mockClient.Client,
	}

	// Define the observable for the Connect method
	uri := "mongodb://localhost:27017"
	connectObservable := service.Connect(uri)

	// Test the observable's output
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	itemChan := connectObservable.Observe()

	select {
	case item := <-itemChan:
		if item.E != nil {
			t.Errorf("Expected no error, got %v", item.E)
		} else {
			if service.GetClient() == nil {
				t.Errorf("Expected client to be initialized, got nil")
			}
		}
	case <-ctx.Done():
		t.Error("Test timed out while waiting for connection")
	}
}

// TestMongoDBServiceDisconnect tests the Disconnect method of MongoDBService
func TestMongoDBServiceDisconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock MongoDB client
	mockClient := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	// Instantiate MongoDBService with mock client
	service := &storage.MongoDBService{
		Client: mockClient.Client,
	}

	// Define the observable for the Disconnect method
	disconnectObservable := service.Disconnect()

	// Test the observable's output
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	itemChan := disconnectObservable.Observe()

	select {
	case item := <-itemChan:
		if item.E != nil {
			t.Errorf("Expected no error, got %v", item.E)
		} else {
			if service.GetClient() != nil {
				t.Errorf("Expected client to be nil after disconnect, got %v", service.GetClient())
			}
		}
	case <-ctx.Done():
		t.Error("Test timed out while waiting for disconnection")
	}
}
