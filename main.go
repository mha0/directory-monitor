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

	// for each folder start goroutine
	results := make(chan Result)

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(config.Dirs))

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
			Check(dir, results)
			waitGroup.Done()
		}(dir)
	}

	// wait for goroutines
	go func() {
		waitGroup.Wait()
		close(results)
	}()

	// write output
	for result := range results {
		log.Println(result.message)
	}
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
