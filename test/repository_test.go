package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"urlshortener/internal/domain"
	"urlshortener/internal/repository"
)

// MockSingleResult simulates a single search result
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// MockCollection mocks URLCollectionInterface
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return nil, args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return nil, args.Error(1)
}

// Test for SaveURL method
func TestSaveURL(t *testing.T) {
	mockCollection := new(MockCollection)
	urlService := &repository.URLServiceImpl{UrlCollection: mockCollection}
	testURL := domain.URL{ID: "testID", OriginalURL: "https://example.com"}

	// Set up the expectation
	mockCollection.On("InsertOne", mock.Anything, testURL).Return(nil, nil)

	// Call SaveURL and observe the result
	observable := urlService.SaveURL(testURL)
	item := <-observable.Observe()

	assert.NoError(t, item.E)
	assert.Equal(t, testURL, item.V.(domain.URL))
	mockCollection.AssertExpectations(t)
}

// Test for GetURL method
func TestGetURL(t *testing.T) {
	mockCollection := new(MockCollection)
	urlService := &repository.URLServiceImpl{UrlCollection: mockCollection}

	// Prepare the expected URL
	expectedURL := domain.URL{ID: "testID", OriginalURL: "https://example.com"}

	// Create a mongo.SingleResult with the expected document
	reg := bson.NewRegistry()
	singleResult := mongo.NewSingleResultFromDocument(expectedURL, nil, reg)

	// Simulate FindOne returning the single result
	mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(singleResult)

	observable := urlService.GetURL("testID")
	item := <-observable.Observe()

	assert.NoError(t, item.E)
	assert.Equal(t, expectedURL, item.V.(domain.URL))
	mockCollection.AssertExpectations(t)
}

// Test for UpdateURL method
func TestUpdateURL(t *testing.T) {
	mockCollection := new(MockCollection)
	urlService := &repository.URLServiceImpl{UrlCollection: mockCollection}
	testURL := domain.URL{ID: "testID", OriginalURL: "https://example.com", Enabled: true}

	// Simulate UpdateOne returning a successful result
	mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	observable := urlService.UpdateURL(testURL)
	item := <-observable.Observe()

	assert.NoError(t, item.E)
	assert.Equal(t, testURL, item.V.(domain.URL))
	mockCollection.AssertExpectations(t)
}
