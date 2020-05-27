package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
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

func IdFromUrl(r *http.Request) (uint, error) {
	var val uint

	key, err := KeyFromUrl(r, "id")
	if err != nil {
		return val, err
	}
	ID, err := strconv.ParseUint(key, 10, 32)
	if err != nil {
		return val, errors.New("Error while decoding request body")
	}
	return uint(ID), nil
}

func KeyFromUrl(r *http.Request, urlKey string) (string, error) {
	var val string
	key, ok := r.URL.Query()[urlKey]
	if !ok || len(key[0]) < 1 {
		logrus.Info("Bad key length")
		logrus.Info(r.URL.Query())

		return val, errors.New("Error while decoding request body")
	}

	return key[0], nil
}
