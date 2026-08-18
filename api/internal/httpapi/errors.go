package httpapi

import (
	"encoding/json"
	"net/http"
)

// ErrorBody is the single JSON error envelope used by the API.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail describes a failure in a stable, client-friendly shape.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message, field string) {
	writeJSON(w, status, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message, Field: field},
	})
}
