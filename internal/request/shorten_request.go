package request

// ShortenRequest defines the structure for URL shortening requests
type ShortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}
