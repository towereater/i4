package main

import (
	"aggregator/config"
	"aggregator/dgpr646"
	"math/rand"
	"os"
	"time"
)

func main() {
	//Setup machine with config
	cfg, err := config.ReadConfig("./config.json")
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(3)
	}

	//Elaboration caches
	dgpr646Cache := &dgpr646.Cache{}

	//Create loop
	for {
		//dgpr646.Fetch()
		dgpr646.Discover(cfg, cfg.FileDir, cfg.Targets[0].File, dgpr646Cache)

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
