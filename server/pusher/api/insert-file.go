package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"pusher/config"
	"pusher/db"
)

func InsertFile(w http.ResponseWriter, r *http.Request) {
	// Extraction of the context config
	//cfg := r.Context().Value(config.ContextConfig).(config.Config)

	// Extraction of extra parameters
	hash := r.PathValue(string(config.ContextHash))
	if hash == "" {
		http.Error(w, "Invalid key", http.StatusForbidden)
		return
	}

	// TODO: GET FILE CONTENT
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
	content := buf.String()
	fmt.Println(content)

	buf.Reset()

	// TODO: CHECK HASH AND METADATA

	// Execution of the request
	err = db.InsertFile(r.Context(), "")
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
