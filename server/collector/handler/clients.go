package handler

import (
	"collector/api"
	"net/http"
)

func ClientsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check of the method request
		switch r.Method {
		case "POST":
			api.InsertClient(w, r)
		default:
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
		}
	})
}

func ClientByIdHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check of the method request
		switch r.Method {
		case "GET":
			api.SelectClient(w, r)
		default:
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
		}
	})
}
