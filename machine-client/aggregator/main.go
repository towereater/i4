package main

import (
	"aggregator/config"
	"aggregator/dgpr646"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Get run args
	fmt.Printf("Starting execution, arg params: %v\n", os.Args)
	if len(os.Args) < 2 {
		fmt.Printf("No config file set\n")
		os.Exit(1)
	}
	configPath := os.Args[1]

	// Setup machine config
	fmt.Printf("Loading configuration from %s\n", configPath)
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Error while reading config file: %s\n", err.Error())
		os.Exit(2)
	}

	// Elaboration caches
	dgpr646Cache := &dgpr646.Cache{}

	// Main loop
	for {
		for _, t := range cfg.Targets {
			switch t.Machine {
			case "DGPR646":
				dgpr646.Fetch(cfg, t)
				continue
				dgpr646.Discover(cfg, t, dgpr646Cache)
			}
		}

		//Choose a remote target and do rename + FTP + delete remote file
		/*
			foreach config target ssh-connect to it using name, user, pass
			check folder for files
			ftp all files and remove them from remote
			load local start-stop
			elaborate them
		*/

		//File sent to remote DB
		/*
			request data save to server
			sent data to server
		*/

		// Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}
