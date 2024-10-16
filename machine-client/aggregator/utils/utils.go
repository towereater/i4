package utils

import (
	"aggregator/config"
	"aggregator/model"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SendFile(cfg config.Config, path string, machine string) {
	// Gathering of the data
	timestamp := time.Now().Format(time.DateTime)

	// Construction of the request
	url := "http://" + cfg.Server.Host + cfg.Server.UploadMetadata
	metadataInput := model.InsertMetadataInput{
		Client:    cfg.Client,
		Machine:   machine,
		Timestamp: timestamp,
		Size:      2,
		Extension: "txt",
		FileHash:  2,
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

	// Opening file connection
	f, err := os.Open(path)
	if err != nil {
		println("Error while opening file:", err)
		return
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		println("Error while opening file:", err)
		return
	}
	_, err = io.Copy(part, f)
	err = writer.Close()

	// Construction of the request
	url = "http://" + metadataOutput.Urls.UploadContent

	// Execution of the request
	res, err = executeHttpRequest(http.MethodPost, url, buf)
	if err != nil {
		println("Error while uploading metadata:", err)
		return
	}
	defer res.Body.Close()

	// Response parsing
	var output model.InsertMetadataOutput
	err = json.NewDecoder(res.Body).Decode(&output)
	if err != nil {
		println("Error while converting metadata output:", err)
		return
	}
}

// func connectSsh() {

// }

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
