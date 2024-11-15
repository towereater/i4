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
	WaterLevel struct {
		Label string  `json:"label"`
		Max   float32 `json:"max"`
		Min   float32 `json:"min"`
	} `json:"waterLevel"`
	Job struct {
		Start string `json:"labelStart"`
		End   string `json:"labelEnd"`
	} `json:"jobLog"`
	UserLog struct {
		Login  string `json:"labelStart"`
		Logoff string `json:"labelEnd"`
	} `json:"userLog"`
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

	//Conversion of the json to struct
	var config Config
	err = json.Unmarshal(byteFile, &config)

	return config, err
}
