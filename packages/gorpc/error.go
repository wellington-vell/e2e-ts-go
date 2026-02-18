package gorpc

import (
	"encoding/json"
	"net/http"
)

// writeError serializes an error to the HTTP response. It differentiates between
// HTTPError (which carries a specific status code) and generic errors (which
// default to 500 Internal Server Error). This approach allows handlers to return
// typed errors with appropriate status codes while maintaining a simple error
// interface for business logic.
func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if httpErr, ok := err.(*HTTPError); ok {
		w.WriteHeader(httpErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{"message": httpErr.Message})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
}

// HTTPError represents an error with an associated HTTP status code.
// It implements the error interface, allowing it to be returned from handlers
// alongside business logic errors. The status code determines the HTTP response
// status, while the message is serialized as JSON in the response body.
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTP error with the specified status code and message.
// This is the primary way to return typed errors from handlers that need specific
// HTTP status codes (e.g., 400 for bad input, 404 for not found, 401 for unauthorized).
func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}
