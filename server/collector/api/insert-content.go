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

func InsertFile(w http.ResponseWriter, r *http.Request) {
	// Extraction of extra parameters
	hash, err := strconv.ParseUint(r.PathValue(string(config.ContextHash)), 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}
	hash32 := uint32(hash)

	// Get the file content
	r.ParseMultipartForm(32 << 20)
	f, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()

	size := header.Size
	fmt.Println(size)

	var buf bytes.Buffer
	io.Copy(&buf, f)
	cccc := buf.Bytes()

	// TODO: CHECK HASH AND METADATA

	// Creation of the file content object
	content := model.UploadContent{
		Hash:    hash32,
		Content: cccc,
	}

	// Execution of the request
	err = db.InsertContent(r.Context(), content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Queue server-side file elaboration
	err = queueContent(r.Context(), hash32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func queueContent(ctx context.Context, hash uint32) error {
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

	// Writing hash on queue
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(datetime),
			Value: h,
		},
	)
	if err != nil {
		return err
	}

	return w.Close()
}
