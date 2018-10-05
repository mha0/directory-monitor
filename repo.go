package main

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
)

const storeFileName = "directory-monitor-store.json"

type Store struct {
	Values map[string]int `json:"Values"`
}

func init() {
	storeFileName := getStoreFileName()
	if _, err := os.Stat(storeFileName); os.IsNotExist(err) {
		store := Store{Values: map[string]int{}}
		WriteStoreToFile(store)
		log.Println("Initialized data store file")
	}
}

func openStoreFile() *os.File {
	storeFile, err := os.OpenFile(getStoreFileName(), os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		log.Fatalln("Could not open store file:", err)
	}
	return storeFile
}

func getStoreFileName() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	configDir := usr.HomeDir + "/.go/"
	storeFileName := configDir + storeFileName
	return storeFileName
}

func ReadStoreFromFile() (store Store) {
	storeFile := openStoreFile()
	defer storeFile.Close()
	decoder := json.NewDecoder(storeFile)
	err := decoder.Decode(&store)
	if err != nil {
		log.Fatalln("Could not decode store file:", err)
	}
	return
}

func WriteStoreToFile(store Store) {
	storeFile := openStoreFile()
	defer storeFile.Close()
	encoder := json.NewEncoder(storeFile)
	err := encoder.Encode(store)
	if err != nil {
		log.Fatalln("Could not encode store file:", err)
	}
}
