package service

import (
	"fmt"
	"github.com/mha0/directory-monitor/domain"
	"os"
	"strings"
)

func renderMessage(dir *os.File, status domain.Status, numberOfFilesAdded int) (message string) {
	switch status {
	case domain.OPERATIONAL:
		message = fmt.Sprintf("%v %v: %v file(s) added since last run.", substringAfterLast(dir.Name(), "/"), status, numberOfFilesAdded)
	case domain.INITIALIZED:
		message = fmt.Sprintf("%v %v: file counter initialized.", substringAfterLast(dir.Name(), "/"), status)
	case domain.WARNING:
		message = fmt.Sprintf("%v %v: No files added since last run!", substringAfterLast(dir.Name(), "/"), status)
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
