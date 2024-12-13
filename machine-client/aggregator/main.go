package main

import (
	"aggregator/config"
	"aggregator/dgpr646"
	"aggregator/utils"
	"fmt"
	"math/rand"
	"os"
	"path"
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
			// Output file name
			//CHANGE FILE NAME
			outputPath := path.Join(cfg.FileDir, fmt.Sprintf("output-%s.txt", t.Machine))

			//CHECK IF ANY FILE IS ALREADY THERE AND ELABORATE IT
			//CYCLE THROUGH LOCAL FILES WITH CORRECT NAME

			// Open output file
			output, err := utils.CreateOrReplaceFile(outputPath)
			if err != nil {
				fmt.Printf("Error while opening output file %s: %s\n", outputPath, err.Error())
				continue
			}

			// Download data file from remote host
			input, err := utils.GetFileFromRemote(cfg, t)
			if err != nil {
				fmt.Printf("Error while fetching file from remote %s: %s\n", t.NetIp, err.Error())
				continue
			}

			// Elaborate data file depending on machine type
			switch t.Machine {
			case "DGPR646":
				err = dgpr646.Elaborate(input, output, dgpr646Cache)
			}
			if err != nil {
				fmt.Printf("Error while elaborating file %s: %s\n", input.Name(), err.Error())
				continue
			}

			// Send file to server
			err = utils.SendFile(cfg, output, t.Machine)
			if err != nil {
				fmt.Printf("Error while sending file %s to server: %s\n", output.Name(), err.Error())
				continue
			}
		}

		// Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}
