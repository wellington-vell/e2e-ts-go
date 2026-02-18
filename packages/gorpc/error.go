package gorpc

import (
	"encoding/json"
	"net/http"
)

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

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}
