package main

import (
	"aggregator/model"
	"fmt"
	"strconv"
)

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
