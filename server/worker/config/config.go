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
		Host        string `json:"host"`
		Port        string `json:"port"`
		Timeout     int    `json:"timeout"`
		DBName      string `json:"dbname"`
		Collections struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Content struct {
				Name string `json:"name"`
			} `json:"content"`
			Interval struct {
				Name string `json:"name"`
			} `json:"interval"`
			Gauge struct {
				Name string `json:"name"`
			} `json:"gauge"`
		} `json:"collections"`
	} `json:"db"`
	Queue struct {
		Host    string `json:"host"`
		Port    string `json:"port"`
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
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
