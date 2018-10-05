package main

import (
	"os"
)

const (
	INITIALIZED = iota
	WARNING
	OPERATIONAL
)

type Status int

func (status Status) String() string {
	// ... operator counts how many
	// items in the array (7)
	names := [...]string{
		"INITIALIZED",
		"WARNING",
		"OPERATIONAL"}

	if status < INITIALIZED || status > OPERATIONAL {
		return "Unknown"
	}

	return names[status]
}

type Result struct {
	file            *os.File
	status          Status
	message         string
	lastRunCount    int
	currentRunCount int
}

func Check(dir *os.File, lastRunCount int, results chan<- Result) {
	// count files
	currentRunCount := CountFiles(dir)

	// map to status
	status := mapToStatus(lastRunCount, currentRunCount)

	// render message
	message := RenderMessage(dir, status, currentRunCount-lastRunCount)

	// write output to channel (and writer later)
	results <- Result{dir, status, message, lastRunCount, currentRunCount}
}

func mapToStatus(lastRunCount int, currentRunCount int) Status {
	var status Status
	if lastRunCount < 0 {
		status = INITIALIZED
	} else {
		if currentRunCount > lastRunCount {
			status = OPERATIONAL
		} else {
			status = WARNING
		}
	}
	return status
}
