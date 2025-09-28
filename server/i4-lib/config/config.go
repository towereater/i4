package config

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type BaseConfig struct {
	Server struct {
		Port string `json:"port" envconfig:"SERVER_PORT"`
	} `json:"server"`
	DB    DBConfig    `json:"db"`
	Queue QueueConfig `json:"queue"`
}

type DBConfig struct {
	Host        string `json:"host" envconfig:"DB_HOST"`
	Timeout     int    `json:"timeout" envconfig:"DB_TIMEOUT"`
	DBName      string `json:"dbname"`
	Collections struct {
		Clients  string `json:"clients"`
		Metadata string `json:"metadata"`
		Content  string `json:"content"`
		Gauge    string `json:"gauge"`
		Interval string `json:"interval"`
	} `json:"collections"`
}

type QueueConfig struct {
	Host    string `json:"host"`
	Timeout int    `json:"timeout"`
	Topic   string `json:"topic"`
}

func loadFromFile(path string, config any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(config)
}

func loadFromEnv(config any) error {
	return envconfig.Process("", config)
}

func LoadConfig(path string, config any) error {
	// Read from config file
	err := loadFromFile(path, config)
	if err != nil {
		return err
	}

	// Load the enviromental variables
	err = loadFromEnv(config)
	return err
}
