package handler

import (
	"collector/api"
	"net/http"
)

func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertMetadata(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}
