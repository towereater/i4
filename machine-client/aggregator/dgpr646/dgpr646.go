package dgpr646

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
	"time"
)

func Discover(cfg config.Config, target model.Target, cache *Cache) {
	//Searches for all files with the given pattern
	files, err := fs.Glob(os.DirFS(cfg.FileDir), target.File)
	if err != nil {
		fmt.Printf("Error while searching %s files: %v", target.File, err)
		return
	}

	//Elaborates all files which were found
	for _, f := range files {
		//Open input file connection
		inputPath := path.Join(cfg.FileDir, f)
		inputFile, err := os.Open(inputPath)
		if err != nil {
			fmt.Printf("Error while opening input file %s: %v\n", f, err)
			continue
		}
		defer inputFile.Close()

		//Open output file connection
		outputPath := path.Join(cfg.FileDir, "elab-"+f[0:strings.LastIndex(f, ".")]+".txt")
		outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error while opening or creating output file %s: %v\n", f, err)
			continue
		}
		defer outputFile.Close()

		//Elaborates given file
		err = elaborate(inputFile, outputFile, cache)

		//If an error occurred, the file is renamed to block next elaborations
		if err != nil {
			timestamp := time.Now().Format(time.DateTime)
			os.Rename(inputPath, path.Join(cfg.FileDir, fmt.Sprintf("ERROR %v %v", timestamp, f)))
			os.Remove(outputPath)
		} else {
			utils.SendFile(cfg, outputPath, target.Machine)
		}
	}
}

func elaborate(inputFile *os.File, outputFile *os.File, cache *Cache) error {
	var content *model.DataContent
	var job *model.DataGauge

	//Read file line by line and split each one
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		data := strings.Split(scanner.Text(), ", ")

		// Format and convert data
		switch data[1] {
		case "PRESSURE":
			content = formatPressure(data)
		case "JOBSTART":
			content, job = formatJob(data)
			cache.Job = job
		case "JOBEND":
			content, job = formatJob(data)
		default:
			continue
		}

		// Save data to file
		fmt.Printf("%+v\n", *content)
		jsonByte, err := json.Marshal(*content)
		if err == nil {
			fmt.Printf("Error while converting content: %v\n", err)
			continue
		}
		fmt.Fprintf(outputFile, "%v\n", jsonByte)

		// Addional elaborations for job interval
		if data[1] == "JOBEND" && cache.Job != nil && job != nil && cache.Job.Value == job.Value {
			content = formatJobInterval(*cache.Job, *job)
			jsonByte, _ = json.Marshal(*content)
			fmt.Fprintf(outputFile, "%v\n", jsonByte)
		}
	}

	return nil
}

func formatPressure(data []string) *model.DataContent {
	value, err := strconv.ParseFloat(data[2], 32)
	if err != nil {
		fmt.Printf("Error while working:\nData:%v\nError:%v\n", data, err)
		return nil
	}

	m := model.DataGauge{
		Timestamp: data[0],
		Key:       data[1],
		Value:     float32(value),
	}

	return &model.DataContent{
		Type:    "GAU",
		Content: m,
	}

	/*
		jsonByte, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("Error while working:\nData:%v\nError:%v\n", data, err)
			return nil
		}

		return &model.DataContent{
			Type:    "GAU",
			Content: string(jsonByte),
		}
	*/
}

func formatJob(data []string) (*model.DataContent, *model.DataGauge) {
	m := model.DataGauge{
		Timestamp: data[0],
		Key:       data[1],
		Value:     data[2],
	}

	return &model.DataContent{
		Type:    "GAU",
		Content: m,
	}, &m

	/*
		jsonByte, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("Error while working:\nData:%v\nError:%v\n", data, err)
			return nil, nil
		}

		return &model.DataContent{
			Type:    "GAU",
			Content: string(jsonByte),
		}, &m
	*/
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

	/*
		jsonByte, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("Error while working:\nData:%v\nError:%v\n", jobStart, err)
			return nil
		}

		return &model.DataContent{
			Type:    "INT",
			Content: string(jsonByte),
		}
	*/
}
