package domain

type Pushover struct {
	AppToken  string `json:"appToken"`
	UserToken string `json:"userToken"`
}

type DirectoryMonitorConfig struct {
	HeartbeatThresholdInHours int `json:"heartbeatThresholdInHours"`
	DeadbeatThresholdInHours  int `json:"deadbeatThresholdInHours"`
	Pushover                  Pushover `json:"pushover"`
	Dirs                      []string `json:"dirs"`
}

