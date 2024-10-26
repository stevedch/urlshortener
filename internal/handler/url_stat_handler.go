package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"urlshortener/internal/service"
)

var URLStatService service.URLStatService

type URLStatHandler struct{}

func NewURLStatHandler() *URLStatHandler {
	return &URLStatHandler{}
}

func (s *URLStatHandler) GetURLStats(c *gin.Context) {
	shortID := c.Param("id")
	statsObservable := URLStatService.GetURLStats(shortID)
	result := <-statsObservable.Observe()
	if result.E != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get URL stats"})
		return
	}
	stats := result.V.(map[string]interface{})
	c.JSON(http.StatusOK, stats)
}
