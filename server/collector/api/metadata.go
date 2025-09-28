package api

import (
	"encoding/json"
	"fmt"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
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

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Create the metadata document
	metadata := model.UploadMetadata(req)

	// Insert the metadata document
	inserted, err := db.UpsertMetadata(cfg.DB, client, metadata)
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
