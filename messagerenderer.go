package main

import (
	"fmt"
	"os"
	"strings"
)

func RenderMessage(dir *os.File, status Status, numberOfFilesAdded int) (message string) {
	switch status {
	case OPERATIONAL:
		message = fmt.Sprintf("%v: %v: %v file(s) added since last run.", status, substringAfterLast(dir.Name(), "/"), numberOfFilesAdded)
	case INITIALIZED:
		message = fmt.Sprintf("%v: %v: Directory file counter initialized.", status, substringAfterLast(dir.Name(), "/"))
	case WARNING:
		message = fmt.Sprintf("%v: %v: No files added since last run!", status, substringAfterLast(dir.Name(), "/"))
	}
	return
}

func substringAfterLast(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
