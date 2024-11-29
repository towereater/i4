package api

import (
	"collector/db"
	"collector/model"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
)

func InsertMetadata(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req model.InsertMetadataInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Compute file id
	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, req.FileHash)

	hash := fnv.New32()
	hash.Write([]byte(fmt.Sprintf("pass%s%s%s%sword", req.Client, req.Machine, req.Timestamp, string(h))))

	// Create the metadata document
	metadata := model.UploadMetadata{
		Client:    req.Client,
		Machine:   req.Machine,
		Timestamp: req.Timestamp,
		Size:      req.Size,
		Extension: req.Extension,
		Hash:      hash.Sum32(),
	}

	// Insert the metadata document
	err = db.InsertMetadata(r.Context(), metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response output
	urls := model.InsertMetadataOutputUrls{
		UploadContent: fmt.Sprintf("%s/uploads/content/%d", r.Host, metadata.Hash),
	}
	output := model.InsertMetadataOutput{
		Id:   metadata.Hash,
		Urls: urls,
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
