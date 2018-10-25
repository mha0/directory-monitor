package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type pushoverRequest struct {
	AppToken  string `json:"token"`
	UserToken string `json:"user"`
	Title     string `json:"title"`
	Message   string `json:"message"`
}

const endpoint = "https://api.pushover.net/1/messages.json"

func SendPushoverNotification(appToken string, userToken string, title string, message string) {

	request := pushoverRequest{appToken, userToken, title, message}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(request)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBodyBytes.Bytes()))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(fmt.Sprintf("Pushover Response Status: %v; Response Body: %v", resp.Status, string(body)))

	if resp.StatusCode != 200 {
		log.Panicln("Failed to send notification using pushover")
	}
}
