package main

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	FilePath string `json:"filePath"`
	WaitTime struct {
		Max float32 `json:"max"`
		Min float32 `json:"min"`
	} `json:"waitTime"`
	Pressure struct {
		Label string  `json:"label"`
		Max   float32 `json:"max"`
		Min   float32 `json:"min"`
	} `json:"pressure"`
	Job struct {
		Start string `json:"labelStart"`
		End   string `json:"labelEnd"`
	} `json:"jobLog"`
}

func ReadConfig(path string) (Config, error) {
	//Read entire config file
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	byteFile, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	//Convert the json to struct
	var config Config
	err = json.Unmarshal(byteFile, &config)

	return config, err
}
