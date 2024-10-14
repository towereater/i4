package dgpr646

import (
	"aggregator/model"
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func Discover(dir string, pattern string, cache *Dgpr646Cache) {
	//Searches for all files with the given pattern
	files, err := fs.Glob(os.DirFS(dir), pattern)
	if err != nil {
		fmt.Printf("Error while searching %s files: %v", pattern, err)
		return
	}

	//Elaborates all files which were found
	for _, v := range files {
		//Open input file connection
		inputPath := path.Join(dir, v)
		inputFile, err := os.Open(inputPath)
		if err != nil {
			fmt.Printf("Error while opening input file %s: %v\n", inputPath, err)
			continue
		}
		defer inputFile.Close()

		//Open output file connection
		outputPath := inputPath[0:strings.LastIndex(inputPath, ".")] + ".json"
		outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Printf("Error while opening output file %s: %v\n", outputPath, err)
			continue
		}
		defer outputFile.Close()

		//Elaborates given file
		err = elaborate(inputFile, outputFile, cache)

		//If an error occurred, the file is renamed to block next elaborations
		if err != nil {
			timestamp := time.Now().Format(time.DateTime)
			os.Rename(inputPath, path.Join(dir, fmt.Sprintf("ERROR %v %v", timestamp, inputPath)))
		}
	}
}

func elaborate(inputFile *os.File, outputFile *os.File, cache *Dgpr646Cache) error {
	//Read file line by line and split each one
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ", ")

		if data[1] == "PRESSURE" {
			generatePressure(outputFile, data)
		} else if data[1] == "JOBSTART" {
			m := generateJob(outputFile, data)

			//Setting cache data
			if m != nil {
				cache.Job = m
			}
		} else if data[1] == "JOBEND" {
			m := generateJob(outputFile, data)

			//Checking cache data
			if cache.Job != nil && m != nil && cache.Job.Value == m.Value {
				generateJobInterval(outputFile, *cache.Job, *m)
			}
		} else {
			continue
		}
	}

	return nil
}

func generatePressure(f *os.File, data []string) *model.FloatGauge {
	value, err := strconv.ParseFloat(data[2], 32)
	if err != nil {
		printError(err, data)
		return nil
	}

	m := model.FloatGauge{
		Timestamp: data[0],
		Key:       data[1],
		Value:     float32(value),
	}

	jsonByte, err := json.Marshal(m)
	if err != nil {
		printError(err, data)
		return nil
	}

	err = printData(f, jsonByte)
	if err != nil {
		printError(err, data)
		return nil
	}

	return &m
}

func generateJob(f *os.File, data []string) *model.StringGauge {
	m := model.StringGauge{
		Timestamp: data[0],
		Key:       data[1],
		Value:     data[2],
	}

	jsonByte, err := json.Marshal(m)
	if err != nil {
		printError(err, data)
		return nil
	}

	err = printData(f, jsonByte)
	if err != nil {
		printError(err, data)
		return nil
	}

	return &m
}

func generateJobInterval(f *os.File, jobStart model.StringGauge, jobEnd model.StringGauge) *model.Interval {
	m := model.Interval{
		Start: jobStart.Timestamp,
		End:   jobEnd.Timestamp,
		Key:   jobStart.Key,
	}

	jsonByte, err := json.Marshal(m)
	if err != nil {
		printError(err, jobStart)
		return nil
	}

	err = printData(f, jsonByte)
	if err != nil {
		printError(err, jobStart)
		return nil
	}

	return &m
}

func printData(f *os.File, jsonByte []byte) error {
	_, err := f.WriteString(fmt.Sprintf("%v,\n", string(jsonByte)))

	return err
}

func printError(err error, data interface{}) {
	fmt.Printf("Error while working:\nData:%v\nError:%v", data, err)
}