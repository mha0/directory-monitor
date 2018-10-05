package main

import (
	"os"
)

const (
	INITIALIZED Status = iota
	FAILED
	OPERATIONAL
)

type Status int

func (status Status) String() string {
	// ... operator counts how many
	// items in the array (7)
	names := [...]string{
		"INITIALIZED",
		"FAILED",
		"OPERATIONAL"}

	if status < INITIALIZED || status > OPERATIONAL {
		return "Unknown"
	}

	return names[status]
}

type Result struct {
	file    *os.File
	status  Status
	message string
	// TODO add lastRunCount currentRunCount
}

func Check(dir *os.File, results chan<- Result) {
	// count files
	initRun, lastRunCount, currentRunCount := CountFiles(dir)

	// define success, failure or init run
	status := mapToStatus(initRun, currentRunCount, lastRunCount)

	// render message success, failure or init
	message := RenderMessage(dir, status)

	// write output to channel (and writer later)
	results <- Result{dir, status, message}
}

func mapToStatus(initRun bool, currentRunCount int, lastRunCount int) Status {
	var status Status
	if initRun {
		status = INITIALIZED
	} else {
		if currentRunCount > lastRunCount {
			status = OPERATIONAL
		} else {
			status = FAILED
		}
	}
	return status
}
