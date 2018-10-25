package service

import (
	"github.com/mha0/directory-monitor/domain"
	"github.com/mha0/directory-monitor/notify"
	"github.com/mha0/directory-monitor/util"
	"os"
)

func Check(dir *os.File, lastRunCount int, results chan<- domain.Result) {
	defer func() {
		if p := recover(); p != nil {
			notify.SendPanicNotification(p)
			os.Exit(1)
		}
	}()

	currentRunCount := util.CountFiles(dir)
	status := mapToStatus(lastRunCount, currentRunCount)
	message := renderMessage(dir, status, currentRunCount-lastRunCount)
	results <- domain.Result{dir, status, message, lastRunCount, currentRunCount}
}

func mapToStatus(lastRunCount int, currentRunCount int) domain.Status {
	var status domain.Status
	if lastRunCount < 0 {
		status = domain.INITIALIZED
	} else {
		if currentRunCount > lastRunCount {
			status = domain.OPERATIONAL
		} else {
			status = domain.WARNING
		}
	}
	return status
}
