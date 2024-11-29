package api

import (
	"bytes"
	"collector/config"
	"collector/db"
	"collector/model"
	"collector/utils"
	"io"
	"net/http"
	"strconv"
)

func InsertContent(w http.ResponseWriter, r *http.Request) {
	// Extract extra parameters
	hash, err := strconv.ParseUint(r.PathValue(string(config.ContextHash)), 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hash32 := uint32(hash)

	// Get the file content
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//f, header, err := r.FormFile("file")
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	//size := header.Size
	//fmt.Println(size)

	var buffer bytes.Buffer
	io.Copy(&buffer, f)
	contentBytes := buffer.Bytes()

	// TODO: CHECK HASH AND METADATA
	// Get the metadata document
	metadata, err := db.SelectMetadata(r.Context(), hash32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the content document
	content := model.UploadContent{
		Hash:    hash32,
		Content: contentBytes,
	}

	// Insert the content document
	err = db.InsertContent(r.Context(), content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Queue content elaboration
	err = utils.QueueContent(r.Context(), hash32, metadata.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
