package web

import (
	"encoding/json"
	"net/http"
)

// Message Build JSON message
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond with JSON message
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if val, ok := data["status"]; ok {
		if !val.(bool) {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	json.NewEncoder(w).Encode(data)
}
