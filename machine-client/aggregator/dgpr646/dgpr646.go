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
