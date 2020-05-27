package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RestActions struct {
	Create     func(http.ResponseWriter, *http.Request)
	Put        func(http.ResponseWriter, *http.Request)
	Patch      func(http.ResponseWriter, *http.Request)
	Get        func(http.ResponseWriter, *http.Request)
	Delete     func(http.ResponseWriter, *http.Request)
	Search     func(http.ResponseWriter, *http.Request)
	SearchPost func(http.ResponseWriter, *http.Request)
	GetByID    func(http.ResponseWriter, *http.Request)
}

func CreateRestRoutes(r *mux.Router, name string, actions RestActions) {
	if actions.Create != nil {
		postRoute := fmt.Sprintf("/api/%s", name)
		r.HandleFunc(postRoute, actions.Create).Methods("POST")
	}
	if actions.Get != nil {
		getRoute := fmt.Sprintf("/api/%s", name)
		r.HandleFunc(getRoute, actions.Get).Methods("GET", "OPTIONS")
	}
	if actions.Put != nil {
		putRoute := fmt.Sprintf("/api/%s", name)
		r.HandleFunc(putRoute, actions.Put).Methods("PUT")
	}
	if actions.Patch != nil {
		patchRoute := fmt.Sprintf("/api/%s", name)
		r.HandleFunc(patchRoute, actions.Patch).Methods("PATCH")
	}
	if actions.Delete != nil {
		deleteRoute := fmt.Sprintf("/api/%s", name)
		r.HandleFunc(deleteRoute, actions.Delete).Methods("DELETE")
	}
	if actions.GetByID != nil {
		getByIDRoute := fmt.Sprintf("/api/%s/id", name)
		r.HandleFunc(getByIDRoute, actions.GetByID).Methods("GET", "OPTIONS")
	}
	if actions.Search != nil {
		searchRoute := fmt.Sprintf("/api/%s/search", name)
		r.HandleFunc(searchRoute, actions.Search).Methods("GET", "OPTIONS")
	}
	if actions.SearchPost != nil {
		searchRoute := fmt.Sprintf("/api/%s/search", name)
		r.HandleFunc(searchRoute, actions.Search).Methods("POST", "OPTIONS")
	}
}
