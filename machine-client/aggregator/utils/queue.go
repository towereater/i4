package utils

import (
	"aggregator/config"
	"aggregator/model"
	"bytes"
	"encoding/json"
	"hash/fnv"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SendFile(cfg config.Config, path string, machine string) error {
	// Send file metadata
	metadataOutput, err := uploadMetadata(cfg, path, machine)
	if err != nil {
		return err
	}

	// Send file
	err = uploadMultiform(cfg, path, metadataOutput.Urls.UploadContent)
	return err
}

func uploadMetadata(cfg config.Config, path string, machine string) (model.InsertMetadataOutput, error) {
	// Open file connection
	f, err := os.Open(path)
	if err != nil {
		return model.InsertMetadataOutput{}, err
	}
	defer f.Close()

	// Retrieve file data
	fi, err := f.Stat()
	if err != nil {
		return model.InsertMetadataOutput{}, err
	}

	// Compute request data
	timestamp := time.Now().Format(time.DateTime)
	size := fi.Size()

	hash := fnv.New32()
	_, err = io.Copy(hash, f)
	if err != nil {
		return model.InsertMetadataOutput{}, err
	}

	// Construct the request
	url := "http://" + cfg.Collector.Host + cfg.Collector.UploadMetadata
	metadataInput := model.InsertMetadataInput{
		Client:    cfg.Client,
		Machine:   machine,
		Timestamp: timestamp,
		Size:      size,
		Extension: "txt",
		FileHash:  hash.Sum32(),
	}

	// Execute the request
	res, err := executeHttpRequest(cfg, http.MethodPost, url, metadataInput)
	if err != nil {
		return model.InsertMetadataOutput{}, err
	}
	defer res.Body.Close()

	// Parse of the request
	var metadataOutput model.InsertMetadataOutput
	err = json.NewDecoder(res.Body).Decode(&metadataOutput)

	return metadataOutput, err
}

func uploadMultiform(cfg config.Config, path string, contentUrl string) error {
	// Create multiform file
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return err
	}

	// Open file connection
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write data on multiform file
	_, err = io.Copy(part, f)
	if err != nil {
		return err
	}
	writer.Close()

	// Construct the request
	url := "http://" + contentUrl

	// Execute the request
	_, err = executeHttpFormFile(cfg, http.MethodPost, url, buffer, writer.FormDataContentType())
	return err
}
