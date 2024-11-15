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
	if len(os.Args) < 2 {
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

	//Elaboration caches
	dgpr646Cache := &dgpr646.Cache{}

	//Main loop
	for {
		//dgpr646.Fetch()
		dgpr646.Discover(cfg, cfg.Targets[0], dgpr646Cache)

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

		//Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}
