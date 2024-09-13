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

var AppConfig Config

func ReadConfig(path string) error {
	//Read entire config file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	byteFile, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	//Conversion of the json to struct
	err = json.Unmarshal(byteFile, &AppConfig)
	if err != nil {
		return err
	}

	return nil
}
