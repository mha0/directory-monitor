package main

import (
	"io/ioutil"
	"os"
)

func CountFiles(dir *os.File) (currentRunCount int) {
	files, _ := ioutil.ReadDir(dir.Name())
	return (len(files))
}
