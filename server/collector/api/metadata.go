package api

import (
	"encoding/json"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"i4-lib/service"
	"net/http"
)

func InsertMetadata(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	client := r.PathValue(string(config.ContextClientCode))
	if client == "" {
		service.Log("Client parameter invalid: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the request
	var req model.InsertMetadataInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		service.Log("Could not convert request body")
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
		service.Log("Error while inserting metadata: %s", err.Error())
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
