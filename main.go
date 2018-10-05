package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

func main() {
	config := ReadConfig()
	fmt.Printf("Checking the following dirs for changes: %v\n", config.Dirs)

	store := ReadStoreFromFile()

	// for each folder start goroutine
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(config.Dirs))

	results := make(chan Result)
	for _, dir := range config.Dirs {
		// check if configured dir is in fact a directory
		isADir := isADir(dir)
		if !isADir {
			log.Fatalln(fmt.Sprintf("File %v is not a directory", dir))
		}

		// open dir
		dir, err := os.Open(dir)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Directory %v cannot be opened", dir))
		}

		// process dir in goroutine
		go func(dir *os.File) {
			Check(dir, getLastRunCount(store, dir), results)
			waitGroup.Done()
		}(dir)
	}

	// wait for goroutines to finish
	go func() {
		waitGroup.Wait()
		close(results)
	}()

	// write results to file
	for result := range results {
		store.Values[result.file.Name()] = result.currentRunCount
		log.Println(result.message)
	}
	WriteStoreToFile(store)
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
