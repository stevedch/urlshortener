package interfaces

import (
	"github.com/gin-gonic/gin"
)

// URLShortenerServiceInterface define los métodos para el servicio de acortador de URLs
type URLShortenerServiceInterface interface {
	ShortenURLHandler(c *gin.Context)
	RedirectURLHandler(c *gin.Context)
	ToggleURLStateHandler(c *gin.Context)
}
