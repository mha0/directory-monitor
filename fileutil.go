package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var FilePath string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Panicln(err)
	}
	FilePath = usr.HomeDir + "/.go/"
}

func OpenFile(filename string) *os.File {
	fileName := FilePath + filename
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Panicln("Could not open file at " + fileName)
	}
	return file
}

func CountFiles(dir *os.File) (currentRunCount int) {
	files, _ := ioutil.ReadDir(dir.Name())
	return (len(files))
}
