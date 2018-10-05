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
		writeStoreToFile(store)
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

func ReadKey(key string) int {
	// load store from file
	store := readStoreFromFile()

	// find value
	if value, exists := store.Values[key]; !exists {
		return -1
	} else {
		return value
	}
}

func readStoreFromFile() (store Store) {
	storeFile := openStoreFile()
	decoder := json.NewDecoder(storeFile)
	err := decoder.Decode(&store)
	if err != nil {
		log.Fatalln("Could not decode store file:", err)
	}
	return
}

func WriteKey(key string, value int) {
	store := readStoreFromFile()
	store.Values[key] = value
	writeStoreToFile(store)
}

func writeStoreToFile(store Store) {
	storeFile := openStoreFile()
	defer storeFile.Close()
	encoder := json.NewEncoder(storeFile)
	err := encoder.Encode(store)
	if err != nil {
		log.Fatalln("Could not encode store file:", err)
	}
}
