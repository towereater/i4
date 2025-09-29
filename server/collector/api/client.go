package api

import (
	"encoding/json"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"i4-lib/service"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func SelectClient(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	code := r.PathValue(string(config.ContextClientCode))
	if code == "" {
		service.Log("Request invalid: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Get the client document
	client, err := db.SelectClientByCode(cfg.DB, code)
	if err == mongo.ErrNoDocuments {
		service.Log("No client with code %s", client.Code)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		service.Log("Error while searching client with code %s: %s", code, err.Error())
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
		service.Log("Request invalid: %+v", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Code == "" || req.Name == "" || req.ApiKey == "" {
		service.Log("Request invalid: %+v", req)
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
		service.Log("Client with code %s already exists", client.Code)
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		service.Log("Error while inserting client %+v: %s", client, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setup client collections and indexes
	err = db.SetupClientCollections(cfg.DB, client.Code)
	if err != nil {
		service.Log("Error while creating collections for client %s: %s", client.Code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusCreated)
}
