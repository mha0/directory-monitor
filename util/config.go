package util

import (
	"encoding/json"
	"github.com/mha0/directory-monitor/domain"
	"log"
	"os"
)

const configFileName = "directory-monitor-conf.json"

func ReadConfig() (config domain.DirectoryMonitorConfig) {
	configFile := OpenFile(configFileName)
	defer configFile.Close()
	config = decodeConfigFile(configFile)
	if len(config.Dirs) == 0 {
		log.Panicln("No Dirs to monitor configured!")
	}
	return
}

func decodeConfigFile(configFile *os.File) (config domain.DirectoryMonitorConfig) {
	decoder := json.NewDecoder(configFile)
	config = domain.DirectoryMonitorConfig{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Panicln("Could not decode Config file:", err)
	}
	return
}
