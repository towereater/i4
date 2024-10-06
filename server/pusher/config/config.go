package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	DB struct {
		Host    string `json:"host"`
		Port    string `json:"port"`
		Timeout int    `json:"timeout"`
	} `json:"db"`
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
