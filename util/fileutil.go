package util

import (
	"fmt"
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

func IsADir(dir string) (isADir bool) {
	info, err := os.Stat(dir)
	if err != nil {
		log.Panicln(fmt.Sprintf("File stat on %v failed", dir))
	}
	if !info.Mode().IsDir() {
		isADir = true
	}
	return true
}