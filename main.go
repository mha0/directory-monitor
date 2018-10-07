package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"sync"
)

func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

var config DirectoryMonitorConfig

func main() {
	defer func() {
		if p := recover(); p != nil {
			SendPanicNotification(p)
			os.Exit(1)
		}
	}()

	if len(os.Args) > 1 {
		FilePath = os.Args[1]
		if !isADir(FilePath) {
			log.Panicln("Argument FilePath is not a directory!")
		}
	} else {
		usr, err := user.Current()
		if err != nil {
			log.Panicln(err)
		}
		FilePath = usr.HomeDir + "/.go/"
	}
	log.Println(fmt.Sprintf("FilePath set to %v", FilePath))

	config = ReadConfig()
	log.Println(fmt.Sprintf("Checking the following dirs for changes: %v", config.Dirs))

	CreateStoreIfNotExists()
	store := ReadStoreFromFile()

	// for each folder start goroutine
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(config.Dirs))

	resultsChannel := make(chan Result)
	for _, dir := range config.Dirs {
		// check if configured dir is in fact a directory
		isADir := isADir(dir)
		if !isADir {
			log.Panicln(fmt.Sprintf("File %v is not a directory", dir))
		}

		// open dir
		dir, err := os.Open(dir)
		defer dir.Close()
		if err != nil {
			log.Panicln(fmt.Sprintf("Directory %v cannot be opened", dir))
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
		log.Panicln(fmt.Sprintf("File stat on %v failed", dir))
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

func SendPanicNotification(p interface{}) {
	messageTitle := "Directory Monitor Panic!"
	messageContent := fmt.Sprintf("panic: %v", p)
	Notify(config.Pushover.AppToken, config.Pushover.UserToken, messageTitle, messageContent)
}
