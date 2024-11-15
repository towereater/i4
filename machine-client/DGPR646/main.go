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
	// Get run args
	if len(os.Args) < 2 {
		println("No config file set")
		os.Exit(1)
	}
	configPath := os.Args[1]

	// Setup machine config
	fmt.Println("Loading configuration")
	cfg, err := ReadConfig(configPath)
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(2)
	}

	//Setup machine with locals
	jobStart = false

	//Create loop
	for {
		if jobStart && rand.Intn(2) == 0 {
			//Randomly generate pressure data
			err = generatePressure(cfg)
			if err != nil {
				fmt.Println("Error while generating pressure data:", err)
				os.Exit(4)
			}
		} else {
			//Randomly generate job start data if a job is not active
			//Randomly generate job end data if a job is active
			err = generateJob(cfg)
			if err != nil {
				fmt.Println("Error while generating job data:", err)
				os.Exit(4)
			}
		}

		//Randomly generate pressure and job errors

		//Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func generatePressure(cfg Config) error {
	//Generate random data
	r := rand.Float32()
	pres := r*(cfg.Pressure.Max-cfg.Pressure.Min) + cfg.Pressure.Min
	datetime := time.Now().Format(time.DateTime)

	//Open output file
	f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	output := fmt.Sprintf("%s, %s, %v\n", datetime, cfg.Pressure.Label, pres)
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}

func generateJob(cfg Config) error {
	//Generate random data
	if !jobStart {
		jobId = rand.Int31()
	}
	datetime := time.Now().Format(time.DateTime)

	//Job status update
	jobStart = !jobStart

	//Open output file
	f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	var output string
	if jobStart {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, cfg.Job.Start, jobId)
	} else {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, cfg.Job.End, jobId)
	}
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}
