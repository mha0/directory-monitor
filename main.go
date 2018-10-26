package main

import (
	"fmt"
	"github.com/mha0/directory-monitor/domain"
	"github.com/mha0/directory-monitor/notify"
	"github.com/mha0/directory-monitor/service"
	"github.com/mha0/directory-monitor/util"
	"io"
	"log"
	"os"
	"os/user"
	"sync"
	"time"
)

var cfg domain.DirectoryMonitorConfig

func main() {
	// handles panic attacks
	defer func() {
		if p := recover(); p != nil {
			notify.SendPanicNotification(p)
			os.Exit(1)
		}
	}()

	// process args
	if len(os.Args) > 1 {
		util.FilePath = os.Args[1]
		if !util.IsADir(util.FilePath) {
			log.Panicln("Argument FilePath is not a directory!")
		}
	} else {
		usr, err := user.Current()
		if err != nil {
			log.Panicln(err)
		}
		util.FilePath = usr.HomeDir + "/.go/"
	}

	// configure logger
	f, err := os.OpenFile(util.FilePath + "directory-monitor-log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	// read config
	cfg = util.ReadConfig()
	log.Println(fmt.Sprintf("Checking the following dirs for changes: %v", cfg.Dirs))

	// initialize and load store
	service.CreateStoreIfNotExists()
	store := service.ReadStore()

	// for each folder start goroutine
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(cfg.Dirs))

	resultsChannel := make(chan domain.Result)
	for _, dir := range cfg.Dirs {
		// check if configured dir is in fact a directory
		isADir := util.IsADir(dir)
		if !isADir {
			log.Panicln(fmt.Sprintf("File %v is not a directory", dir))
		}

		// open dir
		dir, err := os.Open(dir)
		defer dir.Close()
		if err != nil {
			log.Panicln(fmt.Sprintf("Directory %v cannot be opened", dir))
		}

		// process each dir in goroutine
		go func(dir *os.File) {
			service.Check(dir, getLastRunCount(store, dir), resultsChannel)
			waitGroup.Done()
		}(dir)
	}

	// wait for goroutines to finish
	go func() {
		waitGroup.Wait()
		close(resultsChannel)
	}()

	// update store and write file
	var results []domain.Result
	for result := range resultsChannel {
		results = append(results, result)
		store.FileCounters[result.Dir.Name()] = result.CurrentRunCount
		log.Println(result.Message)
	}
	// send push notification if applicable
	thisRunStatus := findHighestSeverityStatus(results)

	shouldNotify, messageTitle := service.ShouldNotify(store.LastRunStatus, thisRunStatus, store.LastNotificationTime, cfg)
	if shouldNotify {
		log.Printf("Sending notification '%v'", messageTitle)
		notify.SendNotification(results, messageTitle)
		store.LastNotificationTime = time.Now()
	}

	store.LastRunStatus = thisRunStatus

	service.WriteStore(store)
}

func getLastRunCount(store domain.Store, dir *os.File) (lastRunCount int) {
	if value, exists := store.FileCounters[dir.Name()]; !exists {
		return -1
	} else {
		return value
	}
}

func findHighestSeverityStatus(results []domain.Result) domain.Status {
	status := domain.OPERATIONAL
	for _, v := range results {
		if v.Status > status {
			status = v.Status
		}
	}
	return status
}