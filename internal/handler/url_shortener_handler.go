package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"urlshortener/internal/interfaces"
	"urlshortener/internal/models"
	"urlshortener/internal/request"
	"urlshortener/internal/service"
)

type URLShortenerHandler struct {
	StatService interfaces.URLStatService
}

// NewURLShortenerHandler creates a new instance of URLShortenerService
func NewURLShortenerHandler() *URLShortenerHandler {
	return &URLShortenerHandler{}
}

func (s *URLShortenerHandler) ShortenURLHandler(c *gin.Context) {
	var req request.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	observable := rxgo.Just(req.OriginalURL)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Calls the URL shortening service
			shortURL, err := service.CreateShortURL(item.(string))
			return shortURL, err
		})

	result := <-observable.Observe()
	if result.E != nil {
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

func (s *URLShortenerHandler) RedirectURLHandler(c *gin.Context) {
	id := c.Param("id")

	observable := rxgo.Just(id)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Calls the service to resolve the original URL
			originalURL, err := service.ResolveURL(item.(string))
			return originalURL, err
		})

	result := <-observable.Observe()
	if result.E != nil || result.V == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Logs the access in statistics
	recordObservable := URLStatService.RecordAccess(id)
	recordResult := <-recordObservable.Observe()
	if recordResult.E != nil {
		// You can add logs here if desired
	}

	// Redirects to the original URL
	c.Redirect(http.StatusFound, result.V.(string))
}

func (s *URLShortenerHandler) ToggleURLStateHandler(c *gin.Context) {
	id := c.Param("id")
	observable := rxgo.Just(id)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Calls the service to toggle the URL state
			updated, err := service.ToggleURLState(item.(string))
			return updated, err
		})
	result := <-observable.Observe()
	if result.E != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL state"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": result.V})
}
