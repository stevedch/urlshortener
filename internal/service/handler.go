package service

import (
	"context"
	"errors"
	"net/http"
	"urlshortener/pkg/models"
	"urlshortener/pkg/request"

	"github.com/gin-gonic/gin"
	"github.com/reactivex/rxgo/v2"
)

// ShortenURLHandler handles requests for URL shortening
func ShortenURLHandler(c *gin.Context) {
	var req request.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	observable := rxgo.Just(req.OriginalURL)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Call the URL shortening service
			shortURL, err := CreateShortURL(item.(string))
			return shortURL, err
		})

	result := <-observable.Observe()
	if result.E != nil {
		// Check if the error is of type APIError
		var apiErr *models.APIError
		if errors.As(result.E, &apiErr) {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": result.V})
}

// RedirectURLHandler handles the redirection of short URLs to their original URLs
func RedirectURLHandler(c *gin.Context) {
	id := c.Param("id")

	observable := rxgo.Just(id)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Call the service to resolve the original URL
			originalURL, err := ResolveURL(item.(string))
			return originalURL, err
		})

	result := <-observable.Observe()
	if result.E != nil || result.V == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Redirect to the original URL
	c.Redirect(http.StatusFound, result.V.(string))
}

// ToggleURLStateHandler handles enabling or disabling short URLs
func ToggleURLStateHandler(c *gin.Context) {
	id := c.Param("id")

	observable := rxgo.Just(id)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Call the service to toggle the URL state
			updated, err := ToggleURLState(item.(string))
			return updated, err
		})

	result := <-observable.Observe()
	if result.E != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": result.V})
}
