package main

import (
	"os"
)

type Status int

const (
	OPERATIONAL Status = iota
	INITIALIZED
	WARNING
)

func (status Status) String() string {
	names := [...]string{
		"OPERATIONAL",
		"INITIALIZED",
		"WARNING"}

	if status < OPERATIONAL || status > WARNING {
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
	defer func() {
		if p := recover(); p != nil {
			SendPanicNotification(p)
			os.Exit(1)
		}
	}()

	currentRunCount := CountFiles(dir)
	status := mapToStatus(lastRunCount, currentRunCount)
	message := RenderMessage(dir, status, currentRunCount-lastRunCount)
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
