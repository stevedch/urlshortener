package interfaces

import (
	"github.com/gin-gonic/gin"
)

// URLShortenerServiceInterface defines the methods for the URL shortener service
type URLShortenerServiceInterface interface {
	ShortenURLHandler(c *gin.Context)
	RedirectURLHandler(c *gin.Context)
	ToggleURLStateHandler(c *gin.Context)
}
