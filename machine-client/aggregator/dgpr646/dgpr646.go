package dgpr646

import (
	"aggregator/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
func Elaborate(cfg config.Config, target model.Target, cache *Cache) {
	// Search for all files with the given pattern
	files, err := fs.Glob(os.DirFS(cfg.FileDir), target.File)
	if err != nil {
		fmt.Printf("Error while searching files with pattern %s: %s\n", target.File, err.Error())
		return
	}

	// Elaborate all found files
	for _, f := range files {
		// Open input file connection
		inputPath := path.Join(cfg.FileDir, f)
		inputFile, err := os.Open(inputPath)
		if err != nil {
			fmt.Printf("Error while opening input file %s: %s\n", inputPath, err.Error())
			continue
		}
		defer inputFile.Close()

		// Open output file connection
		outputPath := path.Join(cfg.FileDir, "elab-"+f[0:strings.LastIndex(f, ".")]+".txt")
		outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error while opening output file %s: %s\n", outputPath, err.Error())
			continue
		}
		defer outputFile.Close()

		// Elaborate given file
		err = elaborate(inputFile, outputFile, cache)

		// Send data to queue if no error occurred
		// Rename file to block next elaborations if an error occurred
		if err == nil {
			err = utils.SendFile(cfg, outputPath, target.Machine)
			if err != nil {
				fmt.Printf("Error while sending file %s to server: %s\n", outputPath, err.Error())
				continue
			}
		} else {
			fmt.Printf("Error while elaborating file %s: %s\n", outputPath, err.Error())

			timestamp := time.Now().Format(time.DateTime)
			os.Rename(inputPath, path.Join(cfg.FileDir, fmt.Sprintf("error-%s-%s", timestamp, f)))
			os.Remove(outputPath)
		}
	}
}
*/

func Elaborate(input *os.File, output *os.File, cache *Cache) error {
	// Read file lines and split them
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// Addional elaboration data
		var content *model.DataContent
		var job *model.DataGauge
		var err error

		data := strings.Split(scanner.Text(), ", ")

		// Format and convert data
		switch data[1] {
		case "PRESSURE":
			content, err = formatPressure(data)
		case "JOBSTART":
			content, job, err = formatJob(data)
			if err == nil {
				cache.Job = job
			}
		case "JOBEND":
			content, job, err = formatJob(data)
		default:
			return fmt.Errorf("unrecognized data record: %s", data[1])
		}
		if err != nil {
			return err
		}

		// Save data to file
		jsonByte, err := json.Marshal(content)
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%s\n", string(jsonByte))

		// Elaborate interval if job end is recognized
		if data[1] == "JOBEND" && cache.Job != nil && job != nil && cache.Job.Value == job.Value {
			content = formatJobInterval(*cache.Job, *job)
			jsonByte, err = json.Marshal(content)
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%s\n", string(jsonByte))
		}
	}

	return nil
}

func formatPressure(data []string) (*model.DataContent, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("invalid %s data record", data[1])
	}

	value, err := strconv.ParseFloat(data[2], 32)
	if err != nil {
		return nil, fmt.Errorf("invalid %s data record", data[1])
	}

	m := model.DataGauge{
		Timestamp: data[0],
		Key:       data[1],
		Value:     float32(value),
	}

	return &model.DataContent{
		Type:    "GAU",
		Content: m,
	}, nil
}

func formatJob(data []string) (*model.DataContent, *model.DataGauge, error) {
	if len(data) < 2 {
		return nil, nil, fmt.Errorf("invalid %s data record", data[1])
	}

	m := model.DataGauge{
		Timestamp: data[0],
		Key:       data[1],
		Value:     data[2],
	}

	return &model.DataContent{
		Type:    "GAU",
		Content: m,
	}, &m, nil
}

func formatJobInterval(jobStart model.DataGauge, jobEnd model.DataGauge) *model.DataContent {
	m := model.DataInterval{
		Start: jobStart.Timestamp,
		End:   jobEnd.Timestamp,
		Key:   jobStart.Key,
	}

	return &model.DataContent{
		Type:    "INT",
		Content: m,
	}
}
