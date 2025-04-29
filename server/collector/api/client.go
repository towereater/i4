package api

import (
	"collector/config"
	"collector/db"
	"collector/model"
	"encoding/json"
	"fmt"
	"net/http"
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
	if err != nil {
		fmt.Printf("Error while inserting client %+v: %s\n", client, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusCreated)
}

func SelectClient(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	code := r.PathValue(string(config.ContextClientCode))

	// Get the client document
	client, err := db.SelectClientByCode(r.Context(), code)
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
