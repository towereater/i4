package main

import (
	"dgpr646/config"
	"dgpr646/utils"
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

	// Setup machine with locals
	jobStart = false

	// Main loop
	for {
		if jobStart && rand.Intn(2) == 0 {
			// Generate random pressure data
			err = generatePressure(cfg)
			if err != nil {
				fmt.Printf("Error while generating pressure data: %s\n", err.Error())
			}
		} else {
			// Generate random job start data if a job is not active
			// Generate random job end data if a job is active
			err = generateJob(cfg)
			if err != nil {
				fmt.Printf("Error while generating job data: %s\n", err.Error())
			}
		}

		// Generate random pressure and job errors

		// Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func generatePressure(cfg config.Config) error {
	// Generate random data
	r := rand.Float32()
	pres := r*(cfg.Pressure.Max-cfg.Pressure.Min) + cfg.Pressure.Min
	datetime := time.Now().Format(time.DateTime)

	// Open output file
	f, err := utils.CreateOrAppendFile(cfg.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	output := fmt.Sprintf("%s, %s, %f\n", datetime, cfg.Pressure.Label, pres)
	fmt.Print(output)
	_, err = f.WriteString(output)

	return err
}

func generateJob(cfg config.Config) error {
	// Generate random data
	if !jobStart {
		jobId = rand.Int31()
	}
	datetime := time.Now().Format(time.DateTime)

	// Job status update
	jobStart = !jobStart

	// Open output file
	f, err := utils.CreateOrAppendFile(cfg.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write output file
	var output string
	if jobStart {
		output = fmt.Sprintf("%s, %s, %d\n", datetime, cfg.Job.Start, jobId)
	} else {
		output = fmt.Sprintf("%s, %s, %d\n", datetime, cfg.Job.End, jobId)
	}
	fmt.Print(output)
	_, err = f.WriteString(output)

	return err
}
