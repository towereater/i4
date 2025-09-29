package api

import (
	"encoding/json"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"i4-lib/service"
	"net/http"
)

func InsertMachine(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	client := r.PathValue(string(config.ContextClientCode))
	if client == "" {
		service.Log("Request invalid: %s", r.URL.Path)
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
		service.Log("Request invalid: %+v", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Create the machine document
	machine := model.Machine(req)

	// Insert the machine document
	err = db.InsertMachine(cfg.DB, client, machine)
	if err != nil {
		service.Log("Error while inserting machine %+v: %s", machine, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusCreated)
}

func RemoveMachine(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	client := r.PathValue(string(config.ContextClientCode))
	machine := r.PathValue(string(config.ContextMachineCode))
	if client == "" || machine == "" {
		service.Log("Request invalid: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Remove the machine document
	err := db.RemoveMachine(cfg.DB, client, machine)
	if err != nil {
		service.Log("Error while removing machine %+v: %s", machine, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.WriteHeader(http.StatusNoContent)
}
