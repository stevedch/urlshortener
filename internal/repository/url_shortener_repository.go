package repository

import (
	"context"
	"errors"
	"github.com/reactivex/rxgo/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"urlshortener/internal/domain"
	"urlshortener/internal/interfaces"
)

// URLServiceImpl implements URLServiceInterface
type URLServiceImpl struct {
	UrlCollection interfaces.URLCollectionInterface
}

// InitDatabase initializes the MongoDB connection and assigns the collection
func (s *URLServiceImpl) InitDatabase(client *mongo.Client, dbName, collectionName string) interfaces.URLCollectionInterface {
	s.UrlCollection = client.Database(dbName).Collection(collectionName)
	return s.UrlCollection
}

// SaveURL saves a URL to the database reactively
func (s *URLServiceImpl) SaveURL(url domain.URL) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := s.UrlCollection.InsertOne(ctx, url)
		if err != nil {
			ch <- rxgo.Error(errors.New("failed to save URL"))
		} else {
			ch <- rxgo.Of(url)
		}
	}})
}

// GetURL retrieves a URL from the database by its ID reactively
func (s *URLServiceImpl) GetURL(shortID string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var url domain.URL
		filter := bson.M{"id": shortID}
		err := s.UrlCollection.FindOne(ctx, filter).Decode(&url)
		if errors.Is(err, mongo.ErrNoDocuments) {
			ch <- rxgo.Error(errors.New("URL not found"))
		} else if err != nil {
			ch <- rxgo.Error(err)
		} else {
			ch <- rxgo.Of(url)
		}
	}})
}

// UpdateURL updates the enabled state or original URL in the database reactively
func (s *URLServiceImpl) UpdateURL(url domain.URL) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"id": url.ID}
		update := bson.M{"$set": bson.M{"enabled": url.Enabled, "original_url": url.OriginalURL}}
		_, err := s.UrlCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			ch <- rxgo.Error(errors.New("failed to update URL"))
		} else {
			ch <- rxgo.Of(url)
		}
	}})
}

// FindURLByOriginal searches for a URL by its original URL in the database reactively
func (s *URLServiceImpl) FindURLByOriginal(originalURL string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var url domain.URL
		filter := bson.M{"original_url": originalURL}
		err := s.UrlCollection.FindOne(ctx, filter).Decode(&url)
		if err != nil {
			ch <- rxgo.Error(err) // Returns error if not found
		} else {
			ch <- rxgo.Of(url)
		}
	}})
}
