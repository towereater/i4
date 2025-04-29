package api

import (
	"collector/config"
	"collector/db"
	"collector/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func InsertMachine(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	client := r.PathValue(string(config.ContextClientCode))
	if client == "" {
		fmt.Printf("Request invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the request
	var req model.InsertMachineInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Code == "" || req.Name == "" {
		fmt.Printf("Request invalid: %+v\n", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the machine document
	machine := model.Machine{
		Code: req.Code,
		Name: req.Name,
	}

	// Insert the machine document
	err = db.InsertMachine(r.Context(), client, machine)
	if err != nil {
		fmt.Printf("Error while inserting machine %+v: %s\n", machine, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusCreated)
}

func RemoveMachine(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	client := r.PathValue(string(config.ContextClientCode))
	machine := r.PathValue(string(config.ContextMachineCode))
	if client == "" || machine == "" {
		fmt.Printf("Request invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove the machine document
	err := db.RemoveMachine(r.Context(), client, machine)
	if err != nil {
		fmt.Printf("Error while removing machine %+v: %s\n", machine, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusNoContent)
}
