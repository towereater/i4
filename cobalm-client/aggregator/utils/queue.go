package utils

import (
	"aggregator/config"
	"aggregator/model"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func SendFile(cfg config.Config, filename string) {
	// Open file to get file stats
	filepath := path.Join(cfg.FileDir, filename)
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error while opening file %s: %s\n", filename, err.Error())
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Error while reading file %s stats: %s\n", filename, err.Error())
		return
	}

	// Compute request data
	timestamp := time.Now().Format(time.DateTime)
	size := fi.Size()

	h := sha256.New()
	io.Copy(h, f)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	metadata := model.InsertMetadataInput{
		Timestamp: timestamp,
		Size:      size,
		Extension: "txt",
		Hash:      hash,
	}

	// Send file metadata
	err = uploadMetadata(cfg, metadata)
	if err != nil {
		fmt.Printf("Error while uploading file %s metadata: %s\n", filename, err.Error())
		return
	}

	// Restart file reading position
	f.Seek(0, 0)

	// Send file
	err = uploadMultiform(cfg, f, metadata.Hash)
	if err != nil {
		fmt.Printf("Error while uploading file %s content: %s\n", filename, err.Error())
		return
	}
}

func uploadMetadata(cfg config.Config, metadata model.InsertMetadataInput) error {
	// Construct the request
	url := fmt.Sprintf("http://%s/clients/%s/%s",
		cfg.Collector.Host,
		cfg.Client,
		cfg.Collector.UploadMetadata)

	// Execute the request
	_, err := executeHttpRequest(cfg, http.MethodPost, url, metadata)
	return err
}

func uploadMultiform(cfg config.Config, f *os.File, metadataHash string) error {
	// Create multiform file
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
	if err != nil {
		return err
	}

	// Write data on multiform file
	_, err = io.Copy(part, f)
	if err != nil {
		return err
	}
	writer.Close()

	// Construct the request
	url := fmt.Sprintf("http://%s/clients/%s/%s/%s",
		cfg.Collector.Host,
		cfg.Client,
		cfg.Collector.UploadContent,
		metadataHash)

	// Execute the request
	_, err = executeHttpFormFile(cfg, http.MethodPost, url, buffer, writer.FormDataContentType())
	return err
}
