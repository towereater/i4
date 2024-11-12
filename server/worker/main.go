package main

import (
	"fmt"
	"os"

	"worker/config"
)

func main() {
	// Get run args
	if len((os.Args)) < 2 {
		println("No config file set")
		os.Exit(1)
	}
	configPath := os.Args[1]

	// Setup machine config
	fmt.Println("Loading configuration")
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(2)
	}

	// Starting up
	fmt.Println(cfg)
}
