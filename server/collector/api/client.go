package api

import (
	"encoding/json"
	"fmt"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func SelectClient(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	code := r.PathValue(string(config.ContextClientCode))
	if code == "" {
		fmt.Printf("Request invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Get the client document
	client, err := db.SelectClientByCode(cfg.DB, code)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No client with code %s", client.Code)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Printf("Error while searching client with code %s: %s\n", code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func InsertClient(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req model.InsertClientInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Request invalid: %+v\n", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Code == "" || req.Name == "" || req.ApiKey == "" {
		fmt.Printf("Request invalid: %+v\n", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Create the client document
	client := model.Client(req)

	// Insert the client document
	err = db.InsertClient(cfg.DB, client)
	if mongo.IsDuplicateKeyError(err) {
		fmt.Printf("Client with code %s already exists\n", client.Code)
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		fmt.Printf("Error while inserting client %+v: %s\n", client, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setup client collections and indexes
	err = db.SetupClientCollections(cfg.DB, client.Code)
	if err != nil {
		fmt.Printf("Error while creating collections for client %s: %s\n", client.Code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusCreated)
}
