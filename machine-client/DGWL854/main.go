package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Job status
var jobStart bool
var jobId int32

// User status
var userLogged bool
var userId int32

func main() {
	//Setup machine with config
	err := ReadConfig("./config.json")
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(3)
	}

	//Setup machine with locals
	jobStart = false
	userLogged = false

	//Create loop
	for {
		choice := rand.Intn(4)

		if jobStart && choice > 1 {
			//Randomly generate water level data
			err = generateWaterLevel()
			if err != nil {
				fmt.Println("Error while generating water level data:", err)
				os.Exit(4)
			}
		} else if userLogged && choice == 1 {
			//Randomly generate job start data if a job is not active
			//Randomly generate job end data if a job is active
			err = generateJob()
			if err != nil {
				fmt.Println("Error while generating job data:", err)
				os.Exit(4)
			}
		} else {
			//Randomly generate user login data if a user is logged
			//Randomly generate user logoff data if a user is not logged
			err = generateUserLog()
			if err != nil {
				fmt.Println("Error while generating user log data:", err)
				os.Exit(4)
			}
		}

		//Randomly generate pressure and job errors

		//Wait some time
		r := rand.Float32()
		waitTime := r*(AppConfig.WaitTime.Max-AppConfig.WaitTime.Min) + AppConfig.WaitTime.Min
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func generateWaterLevel() error {
	//Generate random data
	r := rand.Float32()
	water := r*(AppConfig.WaterLevel.Max-AppConfig.WaterLevel.Min) + AppConfig.WaterLevel.Min
	datetime := time.Now().Format(time.DateTime)

	//Open output file
	f, err := os.OpenFile(AppConfig.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	output := fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.WaterLevel.Label, water)
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}

func generateJob() error {
	//Generate random data
	if !jobStart {
		jobId = rand.Int31()
	}
	datetime := time.Now().Format(time.DateTime)

	//Job status update
	jobStart = !jobStart

	//Open output file
	f, err := os.OpenFile(AppConfig.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	var output string
	if jobStart {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.Job.Start, jobId)
	} else {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.Job.End, jobId)
	}
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}

func generateUserLog() error {
	//Generate random data
	if !userLogged {
		userId = rand.Int31()
	}
	datetime := time.Now().Format(time.DateTime)

	//User status update
	userLogged = !userLogged

	//Open output file
	f, err := os.OpenFile(AppConfig.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write output file
	var output string
	if userLogged {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.UserLog.Login, userId)
	} else {
		output = fmt.Sprintf("%s, %s, %v\n", datetime, AppConfig.UserLog.Logoff, userId)
	}
	fmt.Print(output)
	_, err = f.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}
