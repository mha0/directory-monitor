package notify

import (
	"bytes"
	"fmt"
	"github.com/mha0/directory-monitor/domain"
	"github.com/mha0/directory-monitor/util"
)

func SendNotification(results []domain.Result, messageTitle string) {
	messageContent := renderMessageContent(results)
	sendPushNotification(messageTitle, messageContent)
}

func renderMessageContent(results []domain.Result) string {
	var buffer bytes.Buffer
	for _, result := range results {
		buffer.WriteString(" - " + result.Message+ "\n")
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