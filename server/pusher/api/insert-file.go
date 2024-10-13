package api

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"pusher/config"
	"pusher/db"
	"pusher/model"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
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
	err = queueFile(r.Context(), hash32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func queueFile(ctx context.Context, hash uint32) error {
	// Extracting config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Creating topic writer with timeout
	w := &kafka.Writer{
		Addr:  kafka.TCP(fmt.Sprintf("%s:%s", cfg.Queue.Host, cfg.Queue.Port)),
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
