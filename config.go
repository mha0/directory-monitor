package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFileName = "directory-monitor-conf.json"

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

func ReadConfig() (config DirectoryMonitorConfig) {
	configFile := OpenFile(configFileName)
	defer configFile.Close()
	config = decodeConfigFile(configFile)
	if len(config.Dirs) == 0 {
		log.Panicln("No Dirs to monitor configured!")
	}
	return
}

func decodeConfigFile(configFile *os.File) (config DirectoryMonitorConfig) {
	decoder := json.NewDecoder(configFile)
	config = DirectoryMonitorConfig{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Panicln("Could not decode config file:", err)
	}
	return
}
