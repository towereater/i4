package main

import (
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

	//Create loop
	for {
		//Check already locally saved files

		//Choose a remote target and do rename + FTP + delete remote file

		//Events translated immediately

		//Start-stop saved locally

		//File sent to remote DB

		//Setup machine with locals

		//Wait some time
		r := rand.Float32()
		waitTime := r*(AppConfig.WaitTime.Max-AppConfig.WaitTime.Min) + AppConfig.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func connectSsh() {

}
