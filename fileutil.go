package main

import (
	"io/ioutil"
	"log"
	"os"
)

var FilePath string

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
