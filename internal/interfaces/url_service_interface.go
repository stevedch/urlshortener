package interfaces

import (
	"github.com/reactivex/rxgo/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"urlshortener/internal/domain"
)

// URLServiceInterface defines the operations for managing URLs in the database
type URLServiceInterface interface {
	InitDatabase(client *mongo.Client, dbName, collectionName string) URLCollectionInterface
	SaveURL(url domain.URL) rxgo.Observable
	GetURL(shortID string) rxgo.Observable
	UpdateURL(url domain.URL) rxgo.Observable
	FindURLByOriginal(originalURL string) rxgo.Observable
}
