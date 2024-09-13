package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Job status
var jobStart bool
var jobId int32

func main() {
	//Setup machine with config
	err := ReadConfig("./config.json")
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(3)
	}

	//Setup machine with locals
	jobStart = false

	//Create loop
	for {
		if jobStart && rand.Intn(2) == 0 {
			//Randomly generate pressure data
			err = generatePressure()
			if err != nil {
				fmt.Println("Error while generating pressure data:", err)
				os.Exit(4)
			}
		} else {
			//Randomly generate job start data if a job is not active
			//Randomly generate job end data if a job is active
			err = generateJob()
			if err != nil {
				fmt.Println("Error while generating job data:", err)
				os.Exit(4)
			}
		}

		//Randomly generate pressure and job errors

		//Wait some time
		r := rand.Float32()
		waitTime := r*(AppConfig.WaitTime.Max-AppConfig.WaitTime.Min) + AppConfig.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func generatePressure() error {
	//Generate random data
	r := rand.Float32()
	pres := r*(AppConfig.Pressure.Max-AppConfig.Pressure.Min) + AppConfig.Pressure.Min
	datetime := time.DateTime

	//Open output file
	f, err := os.OpenFile(AppConfig.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	output := fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.Pressure.Label, pres)
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}

func generateJob() error {
	//Generate random data
	if !jobStart {
		jobId = rand.Int31()
	}
	datetime := time.DateTime

	//Job status update
	jobStart = !jobStart

	//Open output file
	f, err := os.OpenFile(AppConfig.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	var output string
	if jobStart {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.Job.Start, jobId)
	} else {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.Job.End, jobId)
	}
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}
