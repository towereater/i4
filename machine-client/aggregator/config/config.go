package config

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
	Server struct {
		Host           string `json:"host"`
		UploadMetadata string `json:"uploadMetadata"`
	} `json:"server"`
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
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
