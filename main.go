package main

import (
	"bytes"
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
)

var cfg domain.DirectoryMonitorConfig

func main() {
	defer func() {
		if p := recover(); p != nil {
			notify.SendPanicNotification(p)
			os.Exit(1)
		}
	}()

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

	log.Println(fmt.Sprintf("FilePath set to %v", util.FilePath))
	cfg = util.ReadConfig()
	log.Println(fmt.Sprintf("Checking the following dirs for changes: %v", cfg.Dirs))

	service.CreateStoreIfNotExists()
	store := service.ReadStoreFromFile()

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

		// process dir in goroutine
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

	// write results to file
	results := make(map[string]domain.Result)
	for result := range resultsChannel {
		results[result.Dir.Name()] = result
		store.Values[result.Dir.Name()] = result.CurrentRunCount
		log.Println(result.Message)
	}
	service.WriteStoreToFile(store)

	messageTitle := renderTitle(results)
	messageContent := renderMessageContent(results)
	notify.SendPushNotification(cfg.Pushover.AppToken, cfg.Pushover.UserToken, messageTitle, messageContent)
}

func renderTitle(results map[string]domain.Result) string {
	status := domain.OPERATIONAL
	for _, v := range results {
		if v.Status > status {
			status = v.Status
		}
	}
	return fmt.Sprintf("DirMon Status: %v", status)
}

func renderMessageContent(results map[string]domain.Result) string {
	var buffer bytes.Buffer
	for _, dir := range cfg.Dirs {
		resultMessage := results[dir].Message
		buffer.WriteString(" - " + resultMessage + "\n")
	}
	return buffer.String()
}

func getLastRunCount(store domain.Store, dir *os.File) (lastRunCount int) {
	if value, exists := store.Values[dir.Name()]; !exists {
		return -1
	} else {
		return value
	}
}

