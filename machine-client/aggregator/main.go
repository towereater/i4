package main

import (
	"aggregator/dgpr646"
	"math/rand"
	"os"
	"time"
)

func main() {
	//Setup machine with config
	err := ReadConfig("./config.json")
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(3)
	}

	//Elaboration caches
	var dgpr646Cache *dgpr646.Dgpr646Cache = nil

	//Create loop
	for {
		dgpr646.Discover(AppConfig.FileDir, AppConfig.Targets[0].File, dgpr646Cache)

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
		waitTime := r*(AppConfig.WaitTime.Max-AppConfig.WaitTime.Min) + AppConfig.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func connectSsh() {

}
