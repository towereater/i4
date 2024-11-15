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

func SendFile(cfg config.Config, path string, machine string) {
	// Opening file connection
	f, err := os.Open(path)
	if err != nil {
		println("Error while opening file:", err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		println("Error while getting file size:", err)
		return
	}

	// Gathering of the data
	timestamp := time.Now().Format(time.DateTime)
	size := fi.Size()

	hash := fnv.New32()
	_, err = io.Copy(hash, f)
	if err != nil {
		println("Error while computing file hash:", err)
		return
	}

	// Construction of the request
	url := "http://" + cfg.Collector.Host + cfg.Collector.UploadMetadata
	metadataInput := model.InsertMetadataInput{
		Client:    cfg.Client,
		Machine:   machine,
		Timestamp: timestamp,
		Size:      size,
		Extension: "txt",
		FileHash:  hash.Sum32(),
	}

	// Execution of the request
	res, err := executeHttpRequest(http.MethodPost, url, metadataInput)
	if err != nil {
		println("Error while uploading metadata:", err)
		return
	}
	defer res.Body.Close()

	// Response parsing
	var metadataOutput model.InsertMetadataOutput
	err = json.NewDecoder(res.Body).Decode(&metadataOutput)
	if err != nil {
		println("Error while converting metadata output:", err)
		return
	}
	res.Body.Close()

	// Copying file content to request field
	f.Seek(0, 0)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		println("Error while opening file:", err)
		return
	}

	_, err = io.Copy(part, f)
	if err != nil {
		println("Error while copying file:", err)
		return
	}
	writer.Close()

	// Construction of the request
	url = "http://" + metadataOutput.Urls.UploadContent

	// Execution of the request
	_, err = executeHttpFormFile(http.MethodPost, url, buf, writer.FormDataContentType())
	if err != nil {
		println("Error while uploading metadata:", err)
		return
	}
}

func executeHttpRequest(method string, url string, payload any) (*http.Response, error) {
	// Convertion of the payload
	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Construction of the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bytesPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Execution of the request
	return client.Do(req)
}

func executeHttpFormFile(method string, url string, buf bytes.Buffer, content string) (*http.Response, error) {
	// Construction of the request
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", content)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Execution of the request
	return client.Do(req)
}

// func connectSsh() {

// }
