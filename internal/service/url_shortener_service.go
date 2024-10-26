// internal/service/url_shortener_service.go
package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reactivex/rxgo/v2"

	"urlshortener/internal/interfaces"
	"urlshortener/internal/models"
	"urlshortener/internal/request"
)

type URLShortenerService struct {
	StatService interfaces.URLStatService
}

// NewURLShortenerService crea una nueva instancia de URLShortenerService
func NewURLShortenerService(statService interfaces.URLStatService) *URLShortenerService {
	return &URLShortenerService{
		StatService: statService,
	}
}

func (s *URLShortenerService) ShortenURLHandler(c *gin.Context) {
	var req request.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	observable := rxgo.Just(req.OriginalURL)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Llama al servicio de acortamiento de URL
			shortURL, err := CreateShortURL(item.(string))
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

func (s *URLShortenerService) RedirectURLHandler(c *gin.Context) {
	id := c.Param("id")

	observable := rxgo.Just(id)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Llama al servicio para resolver la URL original
			originalURL, err := ResolveURL(item.(string))
			return originalURL, err
		})

	result := <-observable.Observe()
	if result.E != nil || result.V == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Registra el acceso en las estadísticas
	recordObservable := s.StatService.RecordAccess(id)
	recordResult := <-recordObservable.Observe()
	if recordResult.E != nil {
		// Puedes agregar logs aquí si lo deseas
	}

	// Redirige a la URL original
	c.Redirect(http.StatusFound, result.V.(string))
}

func (s *URLShortenerService) ToggleURLStateHandler(c *gin.Context) {
	id := c.Param("id")
	observable := rxgo.Just(id)().
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Llama al servicio para cambiar el estado de la URL
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
