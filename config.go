package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
)

const configFileName = "directory-monitor.conf"

type DirectoryMonitorConfig struct {
	Dirs []string `json:"dirs"`
}

func (c DirectoryMonitorConfig) String() string {
	return fmt.Sprintf("Dirs: %v", c.Dirs)
}

func ReadConfig() (config DirectoryMonitorConfig) {
	configFile := openConfigFile()
	config = decodeConfigFile(configFile)
	if len(config.Dirs) == 0 {
		log.Fatalln("No Dirs to monitor configured!")
	}
	return
}

func openConfigFile() *os.File {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	configDir := usr.HomeDir + "/.go/"
	configFile := configDir + configFileName
	config, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("Could not read config file at " + configFile)
	}
	return config
}

func decodeConfigFile(configFile *os.File) (config DirectoryMonitorConfig) {
	decoder := json.NewDecoder(configFile)
	config = DirectoryMonitorConfig{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Could not decode config file:", err)
	}
	return
}
