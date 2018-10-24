package notify

import (
	"fmt"
	"github.com/mha0/directory-monitor/util"
)

func SendPanicNotification(p interface{}) {
	messageTitle := "DirMon Panic!"
	messageContent := fmt.Sprintf("panic: %v", p)
	pushover := util.ReadConfig().Pushover
	SendPushNotification(pushover.AppToken, pushover.UserToken, messageTitle, messageContent)
}