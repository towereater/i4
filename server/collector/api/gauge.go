package api

import (
	"encoding/json"
	"fmt"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"i4-lib/service"
	"net/http"
	"strconv"
)

func SelectGauge(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	code := r.PathValue(string(config.ContextClientCode))
	if code == "" {
		service.Log("Request invalid: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	var err error

	machine := queryParams.Get(string(config.ContextMachineCode))
	key := queryParams.Get(string(config.ContextDataKey))
	value := queryParams.Get(string(config.ContextDataValue))
	tsFrom := queryParams.Get(string(config.ContextTimestampFrom))
	tsTo := queryParams.Get(string(config.ContextTimestampTo))

	limit := 50
	if queryParams.Has(string(config.ContextLimit)) {
		limit, err = strconv.Atoi(queryParams.Get(string(config.ContextLimit)))

		if err != nil {
			fmt.Printf("Invalid %s parameter", string(config.ContextLimit))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if limit > 50 {
			limit = 50
		}
	}

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Build the filter
	filter := model.DataGauge{
		Machine: machine,
		Key:     key,
		Value:   value,
	}

	// Select all documents
	data, err := db.SelectGauge(cfg.DB, code, filter, tsFrom, tsTo, limit)
	if err != nil {
		service.Log("Error while searching data for client %s: %s", code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	if len(data) == 0 {
		service.Log("No data found for client %s", code)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func SelectGaugeSum(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	code := r.PathValue(string(config.ContextClientCode))
	if code == "" {
		service.Log("Request invalid: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	var err error

	machine := queryParams.Get(string(config.ContextMachineCode))
	key := queryParams.Get(string(config.ContextDataKey))
	value := queryParams.Get(string(config.ContextDataValue))
	tsFrom := queryParams.Get(string(config.ContextTimestampFrom))
	tsTo := queryParams.Get(string(config.ContextTimestampTo))

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Build the filter
	filter := model.DataGauge{
		Machine: machine,
		Key:     key,
		Value:   value,
	}

	// Select all documents
	data, err := db.SumGauge(cfg.DB, code, filter, tsFrom, tsTo)
	if err != nil {
		service.Log("Error while searching data for client %s: %s", code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	if len(data) == 0 {
		service.Log("No data found for client %s", code)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func SelectGaugeCount(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	code := r.PathValue(string(config.ContextClientCode))
	if code == "" {
		service.Log("Request invalid: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	var err error

	machine := queryParams.Get(string(config.ContextMachineCode))
	key := queryParams.Get(string(config.ContextDataKey))
	value := queryParams.Get(string(config.ContextDataValue))
	tsFrom := queryParams.Get(string(config.ContextTimestampFrom))
	tsTo := queryParams.Get(string(config.ContextTimestampTo))

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Build the filter
	filter := model.DataGauge{
		Machine: machine,
		Key:     key,
		Value:   value,
	}

	// Count all documents
	data, err := db.CountGauge(cfg.DB, code, filter, tsFrom, tsTo)
	if err != nil {
		service.Log("Error while searching data for client %s: %s", code, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
