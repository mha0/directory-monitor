package domain

import "os"

type Result struct {
	Dir             *os.File
	Status          Status
	Message         string
	LastRunCount    int
	CurrentRunCount int
}

