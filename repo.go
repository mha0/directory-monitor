package main

import (
	"encoding/json"
	"log"
	"os"
)

const storeFileName = "directory-monitor-store.json"

type Store struct {
	Values map[string]int `json:"Values"`
}

func init() {
	createStoreIfNotExists()
}

func createStoreIfNotExists() {
	storeFileName := GetDefaultFileLocation() + storeFileName
	if _, err := os.Stat(storeFileName); os.IsNotExist(err) {
		store := Store{Values: map[string]int{}}
		WriteStoreToFile(store)
		log.Println("Initialized data store file")
	}
}

func ReadStoreFromFile() (store Store) {
	storeFile := OpenFile(storeFileName)
	defer storeFile.Close()
	decoder := json.NewDecoder(storeFile)
	err := decoder.Decode(&store)
	if err != nil {
		log.Fatalln("Could not decode store file:", err)
	}
	return
}

func WriteStoreToFile(store Store) {
	storeFile := OpenFile(storeFileName)
	defer storeFile.Close()
	encoder := json.NewEncoder(storeFile)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(store)
	if err != nil {
		log.Fatalln("Could not encode store file:", err)
	}
}
