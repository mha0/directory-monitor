package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
)

func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

var config DirectoryMonitorConfig

func main() {

	// TODO read defaultFileLocation from args

	// TODO use panic/recover to send notification message

	config = ReadConfig()
	fmt.Printf("Checking the following dirs for changes: %v\n", config.Dirs)

	store := ReadStoreFromFile()

	// for each folder start goroutine
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(config.Dirs))

	resultsChannel := make(chan Result)
	for _, dir := range config.Dirs {
		// check if configured dir is in fact a directory
		isADir := isADir(dir)
		if !isADir {
			log.Fatalln(fmt.Sprintf("File %v is not a directory", dir))
		}

		// open dir
		dir, err := os.Open(dir)
		defer dir.Close()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Directory %v cannot be opened", dir))
		}

		// process dir in goroutine
		go func(dir *os.File) {
			Check(dir, getLastRunCount(store, dir), resultsChannel)
			waitGroup.Done()
		}(dir)
	}

	// wait for goroutines to finish
	go func() {
		waitGroup.Wait()
		close(resultsChannel)
	}()

	// write results to file
	results := make(map[string]Result)
	for result := range resultsChannel {
		results[result.file.Name()] = result
		store.Values[result.file.Name()] = result.currentRunCount
		log.Println(result.message)
	}
	WriteStoreToFile(store)

	messageTitle := renderTitle(results)
	messageContent := renderMessageContent(results)
	Notify(config.Pushover.AppToken, config.Pushover.UserToken, messageTitle, messageContent)
}

func renderTitle(results map[string]Result) string {
	status := OPERATIONAL
	for _, v := range results {
		if v.status > status {
			status = v.status
		}
	}
	return fmt.Sprintf("Directory Monitor Status: %v", status)
}

func renderMessageContent(results map[string]Result) string {
	var buffer bytes.Buffer
	for _, dir := range config.Dirs {
		resultMessage := results[dir].message
		buffer.WriteString(" - " + resultMessage + "\n")
	}
	return buffer.String()
}

func isADir(dir string) (isADir bool) {
	info, err := os.Stat(dir)
	if err != nil {
		log.Fatalln(fmt.Sprintf("File stat on %v failed", dir))
	}
	if !info.Mode().IsDir() {
		isADir = true
	}
	return true
}

func getLastRunCount(store Store, dir *os.File) (lastRunCount int) {
	if value, exists := store.Values[dir.Name()]; !exists {
		return -1
	} else {
		return value
	}
}
