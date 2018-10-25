package domain

type Pushover struct {
	UserToken string
	AppToken  string
}

type DirectoryMonitorConfig struct {
	HeartbeatThresholdInHours int
	DeadbeatThresholdInHours  int
	Pushover                  Pushover
	Dirs                      []string
}

