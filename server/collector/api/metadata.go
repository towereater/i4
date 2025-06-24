package api

import (
	"collector/config"
	"collector/db"
	"collector/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func InsertMetadata(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	client := r.PathValue(string(config.ContextClientCode))
	if client == "" {
		fmt.Printf("Client parameter invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the request
	var req model.InsertMetadataInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Could not convert request body\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the metadata document
	metadata := model.UploadMetadata{
		Timestamp: req.Timestamp,
		Size:      req.Size,
		Extension: req.Extension,
		Hash:      req.Hash,
	}

	// Insert the metadata document
	inserted, err := db.UpsertMetadata(r.Context(), client, metadata)
	if err != nil {
		fmt.Printf("Error while inserting metadata: %s\n", err.Error())
		http.Error(w, "Error while inserting metadata", http.StatusInternalServerError)
		return
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	if inserted {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
