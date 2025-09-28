package main

import (
	"aggregator/config"
	"aggregator/model"
	"aggregator/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"os"
	"path"
	"strings"
	"time"
)

// Elaboration constants
const CONTENT_GAUGE = "GAU"
const CONTENT_INTERVAL = "INT"

const COPY_FILENAME = "copy-"
const ELAB_FILENAME = "elab-"

func main() {
	// Get run args
	var configPath string
	log("Starting execution, arg params: %v\n", os.Args)
	if len(os.Args) < 2 {
		configPath = "config.json"
	} else {
		configPath = os.Args[1]
	}

	// Setup machine config
	log("Loading configuration from %s\n", configPath)
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log("Error while reading config file %s: %s\n", configPath, err.Error())
		os.Exit(2)
	}

	// Main loop
	for {
		// Search for all files to elaborate
		files, err := fs.Glob(os.DirFS(cfg.FileDir), cfg.FileRegex)
		if err != nil {
			log("Error while searching files to elaborate: %s\n", err.Error())
			return
		}
		if len(files) > 0 {
			log("Found following files to elaborate: %+v\n", files)
		} else {
			log("Found no files to elaborate\n")
		}

		// Elaborate all found files
		for _, f := range files {
			elaborateFile(cfg, f)
		}

		// Search for all files to send
		files, err = fs.Glob(os.DirFS(cfg.FileDir), ELAB_FILENAME+"*.txt")
		if err != nil {
			log("Error while searching files to send: %s\n", err.Error())
			return
		}
		if len(files) > 0 {
			log("Found following files to send: %+v\n", files)
		} else {
			log("Found no files to send\n")
		}

		// Elaborate all found files
		for _, f := range files {
			filepath := path.Join(cfg.FileDir, f)

			// Send file to server
			err = utils.SendFile(cfg, filepath)
			if err != nil {
				log("Error while sending file %s: %s\n", f, err.Error())
				continue
			}

			// Delete sent file
			if cfg.FileDeletion {
				err = os.Remove(filepath)
				if err != nil {
					log("Error while removing file %s: %s\n", f, err.Error())
					continue
				}
			}
		}

		// Wait some time
		r := rand.Float32()
		waitTime := r*(cfg.WaitTime.Max-cfg.WaitTime.Min) + cfg.WaitTime.Min
		log("Waiting now %f minutes before next elaboration\n", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Minute)
	}
}

func elaborateFile(cfg config.Config, filename string) {
	// Final report data
	var rowRead int32 = 0
	var rowWritten int32 = 0
	var rowSkipped int32 = 0
	var rowError int32 = 0

	// Rename input file
	inputPathOld := path.Join(cfg.FileDir, filename)
	inputPathNew := path.Join(cfg.FileDir, fmt.Sprintf("%s%s", COPY_FILENAME, filename))
	os.Rename(inputPathOld, inputPathNew)

	// Open input file
	inputFile, err := os.Open(inputPathNew)
	if err != nil {
		log("Error while opening input file %s: %s\n", inputPathNew, err.Error())
		return
	}
	defer inputFile.Close()

	// Open output file
	outputPath := path.Join(cfg.FileDir, fmt.Sprintf("%s%s.txt",
		ELAB_FILENAME, filename[0:strings.LastIndex(filename, ".")]))
	outputFile, err := utils.CreateOrAppendFile(outputPath)
	if err != nil {
		log("Error while opening output file %s: %s\n", outputPath, err.Error())
		return
	}
	defer outputFile.Close()

	// Addional elaboration data
	var content model.DataContent
	var job model.DataGauge

	// Elaborate input file
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		// Modify counters
		rowRead++

		// Split data content
		data := strings.Split(scanner.Text(), ", ")

		machine := cfg.Machine

		data[2] = strings.Trim(data[2], "\"")
		data[3] = strings.Trim(data[3], "\"")
		data[4] = strings.Trim(data[4], "\"")

		// Format and convert data
		if data[2] == "INFO" {
			switch data[3] {
			case "LogIN":
				content = getLogOrLoadContent(machine, data)
			case "LogOUT":
				content = getLogOrLoadContent(machine, data)
			case "TIME":
				content = getTimeOrCounterContent(machine, data)
				if content.Type == "" {
					// Modify counters
					rowError++

					continue
				}
			case "COUNTER":
				content = getTimeOrCounterContent(machine, data)
				if content.Type == "" {
					// Modify counters
					rowError++

					continue
				}
			case "PP_Load":
				content = getLogOrLoadContent(machine, data)
			}
		} else if data[2] == "EXECUTION" {
			if data[3] == "PP_Start" {
				job = model.DataGauge{
					Machine:   machine,
					Key:       data[3],
					Value:     data[4],
					Timestamp: getTimestamp(data),
					Params:    nil,
				}
			} else if (data[3] == "PP_End" || data[3] == "PP_Err") && job.Value == data[4] {
				content = model.DataContent{
					Type: CONTENT_INTERVAL,
					Content: model.DataInterval{
						Machine: machine,
						Key:     data[3],
						Value:   data[4],
						Start:   job.Timestamp,
						End:     getTimestamp(data),
						Params:  nil,
					},
				}

				// Save data to file
				jsonByte, err := json.Marshal(&content)
				if err != nil {
					return
				}
				fmt.Fprintf(outputFile, "%s\n", string(jsonByte))

				// Modify counters
				rowWritten++

				// Cache reset
				job = model.DataGauge{}
			}

			content = getLogOrLoadContent(machine, data)
		} else {
			// Modify counters
			rowSkipped++

			log("Unknown record, skipping\n")
			continue
		}

		// Save data to file
		jsonByte, err := json.Marshal(&content)
		if err != nil {
			return
		}
		fmt.Fprintf(outputFile, "%s\n", string(jsonByte))

		// Modify counters
		rowWritten++
	}

	// Print of the final report
	log("Elaboration of file %s complete with:\n"+
		"\tRows read: %d\n"+
		"\tRows written: %d\n"+
		"\tRows skipped: %d\n"+
		"\tRows error: %d\n",
		filename, rowRead, rowWritten, rowSkipped, rowError)

	// Delete input file
	if cfg.FileDeletion {
		inputFile.Close()
		err = os.Remove(inputPathNew)
		if err != nil {
			log("Error while deleting temp file %s: %s\n", inputPathNew, err.Error())
			return
		}
	}
}

func log(message string, params ...any) {
	// Print a message with timestamp data
	fmt.Printf("%s - %s\n",
		time.Now().UTC().Format("2006-01-02T15:04:05"),
		fmt.Sprintf(message, params...),
	)
}
