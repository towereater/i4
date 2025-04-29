package handler

import (
	"collector/api"
	"net/http"
)

func MachinesHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertMachine(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}

func MachinesByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "DELETE":
		api.RemoveMachine(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}
