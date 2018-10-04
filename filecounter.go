package main

import "os"

func CountFiles(dir *os.File) (initRun bool, lastRunCount int, currentRunCount int) {
	return false, 0, 1
}
