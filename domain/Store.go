package domain

import "time"

type Store struct {
	LastRunStatus       Status
	LastTransitionTime  time.Time
	NotificationCounter int
	FileCounters        map[string]int
}
