package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
	} `json:"server"`
	DB struct {
		Host        string `json:"host"`
		Timeout     int    `json:"timeout"`
		DBName      string `json:"dbname"`
		Collections struct {
			Metadata string `json:"metadata"`
			Content  string `json:"content"`
			Gauge    string `json:"gauge"`
			Interval string `json:"interval"`
		} `json:"collections"`
	} `json:"db"`
	Queue struct {
		Brokers []string `json:"brokers"`
		Timeout int      `json:"timeout"`
		Uploads struct {
			Topic string `json:"topic"`
			Group string `json:"group"`
		} `json:"uploads"`
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

	//Convert the json to struct
	var config Config
	err = json.Unmarshal(byteFile, &config)

	return config, err
}
