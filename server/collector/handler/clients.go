package handler

import (
	"collector/api"
	"net/http"
)

func ClientsHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertClient(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}

func ClientsByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "GET":
		api.SelectClient(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}
