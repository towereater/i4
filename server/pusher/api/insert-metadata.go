package api

import (
	"net/http"
	"pusher/db"
)

func InsertMetadata(w http.ResponseWriter, r *http.Request) {
	// Execution of the request
	err := db.InsertMetadata(r.Context(), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
