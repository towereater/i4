package utils

// import (
// 	"aggregator/config"
// 	"aggregator/model"
// 	"bytes"
// 	"hash/fnv"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"time"
// )

// func SendFile(cfg config.Config, filePath string, machine string) error {
// 	// Open input file
// 	f, err := os.Open(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	// Send file metadata
// 	err := uploadMetadata(cfg, f, machine)
// 	if err != nil {
// 		return err
// 	}

// 	// Restart file reading position
// 	f.Seek(0, 0)

// 	// Send file
// 	err = uploadMultiform(cfg, f, metadataOutput.Urls.UploadContent)
// 	return err
// }

// func uploadMetadata(cfg config.Config, f *os.File, machine string) error {
// 	// Retrieve file data
// 	fi, err := f.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	// Compute request data
// 	timestamp := time.Now().Format(time.DateTime)
// 	size := fi.Size()

// 	hash := fnv.New32()
// 	_, err = io.Copy(hash, f)
// 	if err != nil {
// 		return err
// 	}

// 	// Construct the request
// 	url := "http://" + cfg.Collector.Host + cfg.Collector.UploadMetadata
// 	metadataInput := model.InsertMetadataInput{
// 		Timestamp: timestamp,
// 		Size:      size,
// 		Extension: "txt",
// 		FileHash:  hash.Sum32(),
// 	}

// 	// Execute the request
// 	_, err = executeHttpRequest(cfg, http.MethodPost, url, metadataInput)
// 	return err
// }

// func uploadMultiform(cfg config.Config, f *os.File, contentUrl string) error {
// 	// Create multiform file
// 	var buffer bytes.Buffer
// 	writer := multipart.NewWriter(&buffer)
// 	defer writer.Close()

// 	part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
// 	if err != nil {
// 		return err
// 	}

// 	// Write data on multiform file
// 	_, err = io.Copy(part, f)
// 	if err != nil {
// 		return err
// 	}
// 	writer.Close()

// 	// Construct the request
// 	url := "http://" + contentUrl

// 	// Execute the request
// 	_, err = executeHttpFormFile(cfg, http.MethodPost, url, buffer, writer.FormDataContentType())
// 	return err
// }
