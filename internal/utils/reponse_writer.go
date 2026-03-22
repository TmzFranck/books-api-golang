package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseWithError(w http.ResponseWriter, statusCode int, message string, details ...string) {
	response := map[string]any{
		"error": message,
	}
	if len(details) > 0 {
		response["details"] = details[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func ResponseWithData(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
