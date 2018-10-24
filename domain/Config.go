package domain

import "fmt"

type pushover struct {
	UserToken string
	AppToken  string
}

type DirectoryMonitorConfig struct {
	Pushover pushover
	Dirs     []string `json:"dirs"`
}

func (c DirectoryMonitorConfig) String() string {
	return fmt.Sprintf("Dirs: %v", c.Dirs)
}
