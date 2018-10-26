package service

import (
	"fmt"
	"github.com/mha0/directory-monitor/domain"
	"log"
	"time"
)

const transitionTemplate = "Turned %v (from %v)"
const dateFormat = time.RFC822

func ShouldNotify(lastRunStatus domain.Status, thisRunStatus domain.Status, lastNotificationTime time.Time, config domain.DirectoryMonitorConfig) (shouldNotify bool, title string) {
	if thisRunStatus == domain.INITIALIZED {
		return true, fmt.Sprintf("Just %v from scratch", thisRunStatus)
	}

	switch lastRunStatus {
	case domain.INITIALIZED:
		switch thisRunStatus {
		case domain.OPERATIONAL:
			return true, fmt.Sprintf(transitionTemplate, thisRunStatus, lastRunStatus)
		case domain.WARNING:
			return true, fmt.Sprintf(transitionTemplate, thisRunStatus, lastRunStatus)
		}
	case domain.OPERATIONAL:
		switch thisRunStatus {
		case domain.OPERATIONAL:
			if aboveThreshold(lastNotificationTime, config.HeartbeatThresholdInHours) {
				return true, fmt.Sprintf("Heartbeat (still %v)", thisRunStatus)
			} else {
				return false, ""
			}
		case domain.WARNING:
			return true, fmt.Sprintf(transitionTemplate, thisRunStatus, lastRunStatus)
		}
	case domain.WARNING:
		switch thisRunStatus {
		case domain.OPERATIONAL:
			return true, fmt.Sprintf(transitionTemplate, thisRunStatus, lastRunStatus)
		case domain.WARNING:
			if aboveThreshold(lastNotificationTime, config.DeadbeatThresholdInHours) {
				return true, fmt.Sprintf("Deadbeat (still %v)", thisRunStatus)
			} else {
				return false, ""
			}
		}
	}
	log.Panicln(fmt.Sprintf("ShouldNotify case not implemented."))
	return
}

func aboveThreshold(lastNotificationTime time.Time, thresholdInHours int) (aboveThreshold bool) {
	threshold := lastNotificationTime.Add(time.Duration(thresholdInHours) * time.Hour)
	now := time.Now()
	aboveThreshold = now.After(threshold)
	if(aboveThreshold) {
		log.Println(fmt.Sprintf("Above notification threshold: %v (now) > %v (threshold)", now.Format(dateFormat), threshold.Format(dateFormat)))
	} else {
		log.Println(fmt.Sprintf("Below notification threshold: %v (now) < %v (threshold)", now.Format(dateFormat), threshold.Format(dateFormat)))
	}
	return
}
