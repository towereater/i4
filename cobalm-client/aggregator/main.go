package main

import (
	"aggregator/config"
	"aggregator/model"
	"aggregator/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
)

// Elaboration constants
const CONTENT_GAUGE = "GAU"
const CONTENT_INTERVAL = "INT"

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

	// Search for all files to elaborate
	files, err := fs.Glob(os.DirFS(cfg.FileDir), "Y2*M*.csv")
	if err != nil {
		fmt.Printf("Error while searching files to elaborate: %s\n", err.Error())
		return
	}

	// Final report data
	var rowRead int32 = 0
	var rowWritten int32 = 0
	var rowSkipped int32 = 0
	var rowError int32 = 0

	// Elaborate all found files
	for _, f := range files {
		// Open input file
		inputPath := path.Join(cfg.FileDir, f)
		inputFile, err := os.Open(inputPath)
		if err != nil {
			fmt.Printf("Error while opening input file %s: %s\n", inputPath, err.Error())
			continue
		}
		defer inputFile.Close()

		// Open output file
		outputPath := path.Join(cfg.FileDir, fmt.Sprintf("elab-%s", f))
		outputFile, err := utils.CreateFile(outputPath)
		if err != nil {
			fmt.Printf("Error while opening output file %s: %s\n", outputPath, err.Error())
			continue
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

			machine := "TEST_Machine"

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

				fmt.Printf("Unknown record, skipping\n")
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
	}

	// Print of the final report
	fmt.Printf("Execution complete with:\n"+
		"\tRows read: %d\n"+
		"\tRows written: %d\n"+
		"\tRows skipped: %d\n"+
		"\tRows error: %d\n",
		rowRead, rowWritten, rowSkipped, rowError)
}

func getLogOrLoadContent(machine string, data []string) model.DataContent {
	content := model.DataContent{
		Type: CONTENT_GAUGE,
		Content: model.DataGauge{
			Machine:   machine,
			Key:       data[3],
			Value:     data[4],
			Timestamp: getTimestamp(data),
			Params:    nil,
		},
	}

	return content
}

func getTimeOrCounterContent(machine string, data []string) model.DataContent {
	value, err := strconv.ParseFloat(data[6], 64)

	if err != nil {
		fmt.Printf("Error while converting value: %s\n", err.Error())
		return model.DataContent{}
	}

	content := model.DataContent{
		Type: CONTENT_GAUGE,
		Content: model.DataGauge{
			Machine:   machine,
			Key:       data[4],
			Value:     int64(value),
			Timestamp: getTimestamp(data),
			Params:    nil,
		},
	}

	return content
}

func getTimestamp(data []string) string {
	date := data[0]
	time := data[1]

	return fmt.Sprintf("%s-%s-%s %s", date[6:], date[3:5], date[0:2], time)
}
