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

	"go.mongodb.org/mongo-driver/mongo"
)

func SelectData(w http.ResponseWriter, r *http.Request) {
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

	dataType := queryParams.Get(string(config.ContextDataType))
	if dataType != "GAU" && dataType != "INT" {
		fmt.Printf("Invalid %s parameter", string(config.ContextDataType))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	operation := queryParams.Get(string(config.ContextOperation))
	if dataType != "GAU" && dataType != "INT" {
		fmt.Printf("Invalid %s parameter", string(config.ContextOperation))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	// Get all requested data
	var data any

	// Determine which data to get from db
	switch dataType {
	case "GAU":
		switch operation {
		case "ALL":
			filter := model.DataGauge{
				Machine: machine,
				Key:     key,
				Value:   value,
			}
			data, err = db.SelectGauge(cfg.DB, code, filter, tsFrom, tsTo, limit)
		default:
			service.Log("Invalid operation requested: %s", operation)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "INT":
		switch operation {
		case "ALL":
			filter := model.DataInterval{
				Machine: machine,
				Key:     key,
				Value:   value,
			}
			data, err = db.SelectInterval(cfg.DB, code, filter, tsFrom, tsTo, limit)
		default:
			service.Log("Invalid operation requested: %s", operation)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		service.Log("Invalid data type requested: %s", dataType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err == mongo.ErrNoDocuments {
		service.Log("No data with given filters for client %s", code)
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
