package http

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, v any, statusCode int) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err)
	}
}
