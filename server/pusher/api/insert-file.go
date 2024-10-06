package api

import (
	"net/http"
	"pusher/db"
)

func InsertFile(w http.ResponseWriter, r *http.Request) {
	// Extraction of extra parameters
	id := r.PathValue("fileId")
	if id == "" {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Execution of the request
	err := db.InsertFile(r.Context(), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
