package service

import (
	"encoding/json"
	"github.com/mha0/directory-monitor/domain"
	"github.com/mha0/directory-monitor/util"
	"log"
	"os"
)

const storeFileName = "directory-monitor-store.json"

func CreateStoreIfNotExists() {
	storeFileName := util.FilePath + storeFileName
	if _, err := os.Stat(storeFileName); os.IsNotExist(err) {
		store := domain.Store{Values: map[string]int{}}
		WriteStoreToFile(store)
		log.Println("Initialized data store file")
	}
}

func ReadStoreFromFile() (store domain.Store) {
	storeFile := util.OpenFile(storeFileName)
	defer storeFile.Close()
	decoder := json.NewDecoder(storeFile)
	err := decoder.Decode(&store)
	if err != nil {
		log.Panicln("Could not decode store file:", err)
	}
	return
}

func WriteStoreToFile(store domain.Store) {
	storeFile := util.OpenFile(storeFileName)
	defer storeFile.Close()
	encoder := json.NewEncoder(storeFile)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(store)
	if err != nil {
		log.Panicln("Could not encode store file:", err)
	}
}
