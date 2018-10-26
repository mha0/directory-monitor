package notify

import (
	"bytes"
	"fmt"
	"github.com/mha0/directory-monitor/domain"
	"github.com/mha0/directory-monitor/util"
	"sort"
)

func SendNotification(results []domain.Result, messageTitle string) {
	messageContent := renderMessageContent(results)
	sendPushNotification(messageTitle, messageContent)
}

// Sorts the results alphabetically and lists them in a string
func renderMessageContent(results []domain.Result) string {
	var contents []string
	for _, result := range results {
		contents = append(contents, " - "+result.Message+"\n")
	}
	sort.Strings(contents)

	var buffer bytes.Buffer
	for _, content := range contents {
		buffer.WriteString(content)
	}
	return buffer.String()
}

func SendPanicNotification(p interface{}) {
	messageTitle := "DirMon Panic!"
	messageContent := fmt.Sprintf("panic: %v", p)
	sendPushNotification(messageTitle, messageContent)
}

func sendPushNotification(title string, content string) {
	pushover := util.ReadConfig().Pushover
	SendPushoverNotification(pushover.AppToken, pushover.UserToken, title, content)
}
