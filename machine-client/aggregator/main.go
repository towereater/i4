package main

import (
	"aggregator/config"
	"aggregator/dgpr646"
	"aggregator/dgwl854"
	"aggregator/utils"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"strings"
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
	dgwl854Cache := &dgwl854.Cache{}

	// Main loop
	for {
		// Download data files from all remote hosts
		for _, t := range cfg.Targets {
			// Prepare output file name
			timestamp := time.Now().Format(time.DateTime)
			outputPath := path.Join(cfg.FileDir, fmt.Sprintf("dwn-%s-%s-%s.txt", t.Id, t.Machine, timestamp))

			// Download data file from remote host
			err = utils.GetDataFromRemote(cfg, t, outputPath)
			if err != nil {
				fmt.Printf("Error while fetching file from remote %s: %s\n", t.NetIp, err.Error())
				continue
			}
		}

		// Search for all files to elaborate
		files, err := fs.Glob(os.DirFS(cfg.FileDir), "dwn-*.txt")
		if err != nil {
			fmt.Printf("Error while searching files to elaborate: %s\n", err.Error())
			return
		}

		// Elaborate all found files
		for _, f := range files {
			// Get file metadata from file name
			metadata := strings.Split(f, "-")
			if len(metadata) < 4 {
				fmt.Printf("File name %s does not match name pattern\n", f)
				continue
			}
			machine := metadata[2]

			// Open input file
			inputPath := path.Join(cfg.FileDir, f)
			inputFile, err := os.Open(inputPath)
			if err != nil {
				fmt.Printf("Error while opening input file %s: %s\n", inputPath, err.Error())
				continue
			}
			defer inputFile.Close()

			// Open output file
			outputPath := path.Join(cfg.FileDir, fmt.Sprintf("elab-%s", f[4:]))
			outputFile, err := utils.CreateFile(outputPath)
			if err != nil {
				fmt.Printf("Error while opening output file %s: %s\n", outputPath, err.Error())
				continue
			}
			defer outputFile.Close()

			// Elaborate data file depending on machine type
			switch machine {
			case "DGPR646":
				err = dgpr646.Elaborate(inputFile, outputFile, dgpr646Cache)
			case "DGWL854":
				err = dgwl854.Elaborate(inputFile, outputFile, dgwl854Cache)
			default:
				err = fmt.Errorf("undefined machine type")
			}
			if err != nil {
				// Rename broken file
				inputFile.Close()
				os.Rename(inputPath, fmt.Sprintf("err-%s", inputPath))

				// Remove output file
				outputFile.Close()
				os.Remove(outputPath)

				fmt.Printf("Error while elaborating file %s: %s\n", inputPath, err.Error())
				continue
			}
		}

		// Search for all files to send
		files, err = fs.Glob(os.DirFS(cfg.FileDir), "elab-*.txt")
		if err != nil {
			fmt.Printf("Error while searching files to elaborate: %s\n", err.Error())
			return
		}

		// Elaborate all found files
		for _, f := range files {
			// Get file metadata from file name
			metadata := strings.Split(f, "-")
			if len(metadata) < 4 {
				fmt.Printf("File name %s does not match name pattern\n", f)
				continue
			}
			machine := metadata[2]

			// Prepare file path
			filePath := path.Join(cfg.FileDir, f)

			// Send file to server
			err = utils.SendFile(cfg, filePath, machine)
			if err != nil {
				fmt.Printf("Error while sending file %s to server: %s\n", filePath, err.Error())
				continue
			}

			// Remove the sent file
			err = os.Remove(filePath)
			if err != nil {
				fmt.Printf("Error while removing file %s: %s\n", filePath, err.Error())
				continue
			}
		}

		// Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}
