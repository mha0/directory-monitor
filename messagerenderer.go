package main

import (
	"fmt"
	"os"
	"strings"
)

// TODO ham add no of files added
func RenderMessage(dir *os.File, status Status) (message string) {
	switch status {
	case INITIALIZED:
		message = fmt.Sprintf("%v: %v: Directory file counter initialized.", status, substringAfterLast(dir.Name(), "/"))
	case FAILED:
		message = fmt.Sprintf("%v: %v: No files added since last run!", status, substringAfterLast(dir.Name(), "/"))
	case OPERATIONAL:
		message = fmt.Sprintf("%v: %v: Files were added since last run.", status, substringAfterLast(dir.Name(), "/"))
	}
	return
}

func substringAfterLast(value string, a string) string {
	// Get substring substringAfterLast a string.
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
