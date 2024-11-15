package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	DB struct {
		Host        string `json:"host"`
		Timeout     int    `json:"timeout"`
		DBName      string `json:"dbname"`
		Collections struct {
			Metadata string `json:"metadata"`
			Content  string `json:"content"`
		} `json:"collections"`
	} `json:"db"`
	Queue struct {
		Host    string `json:"host"`
		Timeout int    `json:"timeout"`
		Topic   string `json:"topic"`
	} `json:"queue"`
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
