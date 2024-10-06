package api

import (
	"encoding/json"
	"net/http"
	"pusher/db"
	"pusher/model"
)

func InsertMetadata(w http.ResponseWriter, r *http.Request) {
	// Parsing of the request
	var req model.InsertMetadataInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: COMPUTE HASH, IT SHOULD NOT COME FROM CLIENT

	// Generation of the new document
	metadata := model.FileMetadata{
		Client:    req.Client,
		Machine:   req.Machine,
		Timestamp: req.Timestamp,
		Size:      req.Size,
		Extension: req.Extension,
		Hash:      "", //req.Hash
	}

	// Execution of the request
	err = db.InsertMetadata(r.Context(), metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response output
	//TODO: SEND BACK HASH AND/OR API URL
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
