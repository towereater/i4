package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"pusher/config"
	"pusher/db"
	"pusher/model"
	"strconv"
)

func InsertFile(w http.ResponseWriter, r *http.Request) {
	// Extraction of extra parameters
	hash, err := strconv.ParseUint(r.PathValue(string(config.ContextHash)), 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hash32 := uint32(hash)

	// Get the file content
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	f, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	size := header.Size
	fmt.Println(size)

	io.Copy(&buf, f)
	cccc := buf.Bytes()
	data := string(cccc)
	fmt.Println(data)

	//buf.Reset()

	// TODO: CHECK HASH AND METADATA

	// Creation of the file content object
	content := model.FileContent{
		Hash:    hash32,
		Content: cccc,
	}

	// Execution of the request
	err = db.InsertFile(r.Context(), content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Queue server-side file elaboration
	err = queueFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

// TODO: ADD FILE HASH TO QUEUE
func queueFile() error {
	return nil
}
