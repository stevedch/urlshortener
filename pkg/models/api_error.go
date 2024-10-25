package models

// APIError represents a custom error with an HTTP status code and a message
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface so that APIError can be used as an error
func (e *APIError) Error() string {
	return e.Message
}
