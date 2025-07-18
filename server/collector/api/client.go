package api

import (
	"collector/config"
	"collector/db"
	"collector/model"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

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

	// Create the client document
	client := model.Client{
		Code:     req.Code,
		Name:     req.Name,
		ApiKey:   req.ApiKey,
		Machines: req.Machines,
	}

	// Insert the client document
	err = db.InsertClient(r.Context(), client)
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
	db.SetupClientCollections(r.Context(), client.Code)
	if err != nil {
		fmt.Printf("Error while creating collections for client %s: %s\n", client.Code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusCreated)
}

func SelectClient(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	code := r.PathValue(string(config.ContextClientCode))
	if code == "" {
		fmt.Printf("Request invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the client document
	client, err := db.SelectClientByCode(r.Context(), code)
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
