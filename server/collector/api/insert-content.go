package api

import (
	"bytes"
	"collector/config"
	"collector/db"
	"collector/model"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func InsertContent(w http.ResponseWriter, r *http.Request) {
	// Extraction of extra parameters
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

	f, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	size := header.Size
	fmt.Println(size)

	var buffer bytes.Buffer
	io.Copy(&buffer, f)
	contentBytes := buffer.Bytes()

	// TODO: CHECK HASH AND METADATA
	metadata, err := db.SelectMetadata(r.Context(), hash32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Creation of the file content object
	content := model.UploadContent{
		Hash:    hash32,
		Content: contentBytes,
	}

	// Execution of the request
	err = db.InsertContent(r.Context(), content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Queue server-side file elaboration
	err = queueContent(r.Context(), hash32, metadata.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func queueContent(ctx context.Context, hash uint32, client string) error {
	// Extracting config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Creating topic writer with timeout
	w := &kafka.Writer{
		Addr:  kafka.TCP(cfg.Queue.Host),
		Topic: cfg.Queue.Topic,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Queue.Timeout)*time.Second)
	defer cancel()

	// Preparing data for queue
	datetime := time.Now().Format(time.DateTime)
	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, hash)
	value := append(h, client...)

	// Writing hash on queue
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(datetime),
			Value: value,
		},
	)
	if err != nil {
		return err
	}

	return w.Close()
}
