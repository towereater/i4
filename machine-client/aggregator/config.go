package main

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	FileDir  string `json:"fileDir"`
	WaitTime struct {
		Max float32 `json:"max"`
		Min float32 `json:"min"`
	} `json:"waitTime"`
	Targets []struct {
		Name   string `json:"name"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
		Folder string `json:"folder"`
		File   string `json:"file"`
	} `json:"targets"`
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
