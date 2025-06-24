package api

import (
	"bytes"
	"collector/config"
	"collector/db"
	"collector/model"
	"collector/utils"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

func InsertContent(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	client := r.PathValue(string(config.ContextClientCode))
	if client == "" {
		fmt.Printf("Client parameter invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := r.PathValue(string(config.ContextHash))
	if hash == "" {
		fmt.Printf("Hash parameter invalid: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the file content
	f, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("File content invalid: %s\n", err.Error())
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

	// Get the metadata document
	metadata, err := db.SelectMetadata(r.Context(), client, hash)
	if err != nil {
		fmt.Printf("No metadata with given hash: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Match metadata with received file properties
	if metadata.Size != sizeCheck {
		fmt.Printf("File size (%d) does not match metadata (%d)\n",
			sizeCheck,
			metadata.Size)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if metadata.Extension != extCheck {
		fmt.Printf("File extension (%s) does not match metadata (%s)\n",
			extCheck,
			metadata.Extension)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if metadata.Hash != hashCheck {
		fmt.Printf("File hash (%s) does not match metadata (%s)\n",
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

	fmt.Printf("File content is:\n%s\n", string(contentBytes))

	content := model.UploadContent{
		Hash:    metadata.Hash,
		Content: contentBytes,
	}

	// Insert the content document
	err = db.InsertContent(r.Context(), client, content)
	if err != nil {
		fmt.Printf("Error while saving content for client %s: %s\n",
			client,
			err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Queue content elaboration request
	err = utils.QueueContent(r.Context(), client, metadata.Hash)
	if err != nil {
		fmt.Printf("Error while queueing content for client %s: %s\n",
			client,
			err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
