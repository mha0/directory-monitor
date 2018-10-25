package service

import (
	"fmt"
	"github.com/mha0/directory-monitor/domain"
	"log"
	"time"
)

const TRANSITION_TEMPLATE = "Turned %v (from %v)"

func ShouldNotify(lastRunStatus domain.Status, thisRunStatus domain.Status, lastTransitionTime time.Time, notificationCounter int, config *domain.DirectoryMonitorConfig) (shouldNotify bool, title string) {
	if thisRunStatus == domain.INITIALIZED {
		return true, fmt.Sprintf("Just %v from scratch", thisRunStatus)
	}

	switch lastRunStatus {
	case domain.INITIALIZED:
		switch thisRunStatus {
		case domain.OPERATIONAL:
			return true, fmt.Sprintf(TRANSITION_TEMPLATE, thisRunStatus, lastRunStatus)
		case domain.WARNING:
			return true, fmt.Sprintf(TRANSITION_TEMPLATE, thisRunStatus, lastRunStatus)
		}
	case domain.OPERATIONAL:
		switch thisRunStatus {
		case domain.OPERATIONAL:
			if aboveThreshold(lastTransitionTime, config.HeartbeatThresholdInHours, notificationCounter) {
				return true, fmt.Sprintf("Heartbeat (still %v)", thisRunStatus)
			} else {
				return false, ""
			}
		case domain.WARNING:
			return true, fmt.Sprintf(TRANSITION_TEMPLATE, thisRunStatus, lastRunStatus)
		}
	case domain.WARNING:
		switch thisRunStatus {
		case domain.OPERATIONAL:
			return true, fmt.Sprintf(TRANSITION_TEMPLATE, thisRunStatus, lastRunStatus)
		case domain.WARNING:
			if aboveThreshold(lastTransitionTime, config.DeadbeatThresholdInHours, notificationCounter) {
				return true, fmt.Sprintf("Deadbeat (still %v)", thisRunStatus)
			} else {
				return false, ""
			}
		}
	}
	log.Panicln(fmt.Sprintf("ShouldNotify case not implemented."))
	return
}

func aboveThreshold(lastTransitionTime time.Time, thresholdInHours int, notificationCounter int) (aboveThreshold bool) {
	x := (notificationCounter + 1) * thresholdInHours
	y := time.Now().Add(time.Duration(-x) * time.Hour)
	return y.After(lastTransitionTime)
}
