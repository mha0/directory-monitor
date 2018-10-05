package main

import "os"

func CountFiles(dir *os.File) (initRun bool, lastRunCount int, currentRunCount int) {
	// count files
	currentRunCount = 2

	// update store
	lastRunCount = ReadKey(dir.Name())
	WriteKey(dir.Name(), currentRunCount)
	if lastRunCount < 0 {
		return true, lastRunCount, currentRunCount
	} else {
		return false, lastRunCount, currentRunCount
	}
}
