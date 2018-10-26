package domain

import "time"

type Store struct {
	LastRunStatus        Status
	LastNotificationTime time.Time
	FileCounters         map[string]int
}
