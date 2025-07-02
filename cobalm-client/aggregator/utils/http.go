package utils

import (
	"aggregator/config"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func executeHttpRequest(cfg config.Config, method string, url string, payload any) (*http.Response, error) {
	// Convert payload to bytes
	jsonByte, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Construct the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Duration(cfg.Collector.Timeout) * time.Second,
	}

	// Execute the request
	return client.Do(req)
}

func executeHttpFormFile(cfg config.Config, method string, url string, buf bytes.Buffer, content string) (*http.Response, error) {
	// Construct the request
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", content)
	client := &http.Client{
		Timeout: time.Duration(cfg.Collector.Timeout) * time.Second,
	}

	// Execute the request
	return client.Do(req)
}
