package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler returns an HTTP handler that serves the OpenAPI specification as JSON
func Handler(procedures []ProcedureInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spec, err := GenerateSpec(procedures)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate OpenAPI spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(spec); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode OpenAPI spec: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
