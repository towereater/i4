package config

import (
	"aggregator/model"
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	FileDir  string `json:"fileDir"`
	Client   string `json:"client"`
	WaitTime struct {
		Max float32 `json:"max"`
		Min float32 `json:"min"`
	} `json:"waitTime"`
	Targets   []model.Target `json:"targets"`
	Collector struct {
		Host           string `json:"host"`
		Timeout        int    `json:"timeout"`
		UploadMetadata string `json:"uploadMetadata"`
	} `json:"collector"`
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
