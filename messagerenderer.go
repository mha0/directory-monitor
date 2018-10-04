package main

import (
	"fmt"
	"os"
	"strings"
)

func RenderMessage(dir *os.File, status Status) (message string) {
	switch status {
	case INITIALIZED:
		message = fmt.Sprintf("INITIALIZED: %v: Directory file counter initialized.", substringAfterLast(dir.Name(), "/"))
	case FAILED:
		message = fmt.Sprintf("FAILED: %v: No files added since last run!", substringAfterLast(dir.Name(), "/"))
	case OPERATIONAL:
		message = fmt.Sprintf("OPERATIONAL: %v: Files were added since last run.", substringAfterLast(dir.Name(), "/"))
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
