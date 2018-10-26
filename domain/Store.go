package domain

import "time"

type Store struct {
	LastRunStatus        Status `json:"lastRunStatus"`
	LastNotificationTime time.Time `json:"lastNotificationTime"`
	FileCounters         map[string]int `json:"fileCounters"`
}
