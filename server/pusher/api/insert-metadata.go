package api

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"pusher/db"
	"pusher/model"
)

func InsertMetadata(w http.ResponseWriter, r *http.Request) {
	// Parsing of the request
	var req model.InsertMetadataInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Computing file id
	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, req.FileHash)

	hash := fnv.New32()
	hash.Write([]byte("pass"))
	hash.Write([]byte(req.Client))
	hash.Write([]byte(req.Machine))
	hash.Write([]byte(req.Timestamp))
	hash.Write(h)
	hash.Write([]byte("word"))

	// Generation of the metadata document
	metadata := model.FileMetadata{
		Hash:      hash.Sum32(),
		Client:    req.Client,
		Machine:   req.Machine,
		Timestamp: req.Timestamp,
		Size:      req.Size,
		Extension: req.Extension,
	}

	// Execution of the request
	err = db.InsertMetadata(r.Context(), metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creation of the output object
	urls := model.InsertMetadataOutputUrls{
		UploadContent: fmt.Sprintf("%s/%d", r.Host, metadata.Hash),
	}
	output := model.InsertMetadataOutput{
		Id:   metadata.Hash,
		Urls: urls,
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
