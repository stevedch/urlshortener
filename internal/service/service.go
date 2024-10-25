package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"urlshortener/internal/cache"
	"urlshortener/pkg/models"
)

// URLServiceInstance URLService is an instance of the interface URLService which will be injected
var URLServiceInstance URLService

// CreateShortURL generates a shortened URL and stores it in the database and cache
func CreateShortURL(originalURL string) (string, error) {

	// Check if the original URL already exists in the database reactively
	existsObservable := URLServiceInstance.FindURLByOriginal(originalURL)
	existsResult := <-existsObservable.Observe()
	if existsResult.E == nil && existsResult.V.(models.URL).ShortURL != "" {
		return existsResult.V.(models.URL).ShortURL, &models.APIError{
			Code:    http.StatusConflict,
			Message: "URL already exists",
		}
	}

	// Generate a unique ID (hash) for the shortened URL
	shortID := generateShortID(originalURL)
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

	// Create the URL structure
	url := models.URL{
		ID:          shortID,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		Enabled:     true,
	}

	// Save to MongoDB reactively
	saveObservable := URLServiceInstance.SaveURL(url)
	saveResult := <-saveObservable.Observe()
	if saveResult.E != nil {
		return "", saveResult.E
	}

	// Cache in Redis reactively
	cacheObservable := cache.SetURL(shortID, originalURL)
	cacheResult := <-cacheObservable.Observe()
	if cacheResult.E != nil {
		fmt.Printf("Error caching URL in Redis: %v\n", cacheResult.E) // Non-blocking error handling
	}

	return shortURL, nil
}

// ResolveURL retrieves the original URL using the shortened ID
func ResolveURL(shortID string) (string, error) {
	// Try to get the URL from cache reactively
	cacheObservable := cache.GetURL(shortID)
	cacheResult := <-cacheObservable.Observe()
	if cacheResult.E == nil && cacheResult.V.(string) != "" {
		return cacheResult.V.(string), nil
	}

	// If not in cache, search in MongoDB reactively
	dbObservable := URLServiceInstance.GetURL(shortID)
	dbResult := <-dbObservable.Observe()
	if dbResult.E != nil {
		return "", errors.New("URL not found")
	}
	url := dbResult.V.(models.URL)

	// If disabled, return an error
	if !url.Enabled {
		return "", errors.New("URL is disabled")
	}

	// Cache for future requests reactively
	cacheSaveObservable := cache.SetURL(shortID, url.OriginalURL)
	cacheSaveResult := <-cacheSaveObservable.Observe()
	if cacheSaveResult.E != nil {
		fmt.Printf("Error caching URL in Redis: %v\n", cacheSaveResult.E) // Non-blocking error handling
	}

	return url.OriginalURL, nil
}

// ToggleURLState toggles the enabled/disabled state of the URL
func ToggleURLState(shortID string) (bool, error) {
	// Retrieve the URL from the database reactively
	dbObservable := URLServiceInstance.GetURL(shortID)
	dbResult := <-dbObservable.Observe()
	if dbResult.E != nil {
		return false, errors.New("URL not found")
	}
	url := dbResult.V.(models.URL)

	// Toggle the enabled state
	url.Enabled = !url.Enabled

	// Save to MongoDB reactively
	updateObservable := URLServiceInstance.UpdateURL(url)
	updateResult := <-updateObservable.Observe()
	if updateResult.E != nil {
		return false, updateResult.E
	}

	// Remove or update in cache based on the new state
	if url.Enabled {
		cacheUpdateObservable := cache.SetURL(shortID, url.OriginalURL)
		cacheUpdateResult := <-cacheUpdateObservable.Observe()
		if cacheUpdateResult.E != nil {
			fmt.Printf("Error updating cache in Redis: %v\n", cacheUpdateResult.E)
			return false, cacheUpdateResult.E
		}
	} else {
		cacheDeleteObservable := cache.DeleteURL(shortID)
		cacheDeleteResult := <-cacheDeleteObservable.Observe()
		if cacheDeleteResult.E != nil {
			fmt.Printf("Error deleting from cache in Redis: %v\n", cacheDeleteResult.E)
			return false, cacheDeleteResult.E
		}
	}

	return url.Enabled, nil
}

// generateShortID generates a unique ID based on the original URL
func generateShortID(originalURL string) string {
	// Uses hashFunction to generate a hash and takes the first 6 characters
	return hashFunction(originalURL)[:6]
}

// hashFunction generates an MD5 hash of the original URL and converts it to a string
func hashFunction(originalURL string) string {
	hashMd5 := md5.New()
	hashMd5.Write([]byte(originalURL))
	return hex.EncodeToString(hashMd5.Sum(nil))
}
