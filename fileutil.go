package main

import (
	"log"
	"os"
	"os/user"
)

func OpenFile(filename string) *os.File {
	fileDir := GetDefaultFileLocation()
	fileName := fileDir + filename
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln("Could not open file at " + fileName)
	}
	return file
}

func GetDefaultFileLocation() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	return usr.HomeDir + "/.go/"
}
