package service

import (
	"fmt"
	"time"
)

func Log(message string, params ...any) {
	// Print a message with timestamp data
	fmt.Printf("%s - %s\n",
		time.Now().UTC().Format("2006-01-02T15:04:05"),
		fmt.Sprintf(message, params...),
	)
}
