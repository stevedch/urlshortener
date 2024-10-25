package models

// URL represents the structure of a shortened URL in the system
type URL struct {
	ID          string `json:"id" bson:"id"`                     // Unique identifier for the shortened URL
	OriginalURL string `json:"original_url" bson:"original_url"` // The full original URL
	ShortURL    string `json:"short_url" bson:"short_url"`       // The generated shortened URL
	Enabled     bool   `json:"enabled" bson:"enabled"`           // URL status (enabled or disabled)
}
