package api

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/model"
	"i4-lib/service"
	"io"
	"net/http"
	"path/filepath"
)

func InsertContent(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	client := r.PathValue(string(config.ContextClientCode))
	if client == "" {
		service.Log("Client parameter invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := r.PathValue(string(config.ContextHash))
	if hash == "" {
		service.Log("Hash parameter invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the file content
	f, header, err := r.FormFile("file")
	if err != nil {
		service.Log("File content invalid: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()

	// Determine file properties
	sizeCheck := uint32(header.Size)

	extCheck := filepath.Ext(header.Filename)[1:]

	h := sha256.New()
	io.Copy(h, f)
	hashCheck := fmt.Sprintf("%x", h.Sum(nil))

	// Extract config
	cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

	// Get the metadata document
	metadata, err := db.SelectMetadata(cfg.DB, client, hash)
	if err != nil {
		service.Log("No metadata with given hash: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Match metadata with received file properties
	if metadata.Size != sizeCheck {
		service.Log("File size (%d) does not match metadata (%d)\n",
			sizeCheck,
			metadata.Size)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if metadata.Extension != extCheck {
		service.Log("File extension (%s) does not match metadata (%s)\n",
			extCheck,
			metadata.Extension)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if metadata.Hash != hashCheck {
		service.Log("File hash (%s) does not match metadata (%s)\n",
			hashCheck,
			metadata.Hash)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the content document
	var buffer bytes.Buffer
	f.Seek(0, 0)
	io.Copy(&buffer, f)
	contentBytes := buffer.Bytes()

	service.Log("File content is:\n%s\n", string(contentBytes))

	content := model.UploadContent{
		Hash:    metadata.Hash,
		Content: contentBytes,
	}

	// Insert the content document
	err = db.InsertContent(cfg.DB, client, content)
	if err != nil {
		service.Log("Error while saving content for client %s: %s\n",
			client,
			err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Queue content elaboration request
	err = service.QueueContent(cfg.Queue, client, metadata.Hash)
	if err != nil {
		service.Log("Error while queueing content for client %s: %s\n",
			client,
			err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
