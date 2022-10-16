package notifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/template"
)

type SendData struct {
	TrackerUUID string
	TrackerName string
	SubnetCIDR  string
	SubnetName  string
	SubnetTags  []string
	URL         string
	IP          string
	UserAgent   string
	ChatID      []string
}

type Config struct {
	NotificationsURL   string
	NotificationsToken string
}

type request struct {
	ChatID  int64  `json:"chatID"`
	Message string `json:"message"`
}

var config Config
var notificationTemplate string

func init() {
	file, err := os.Open("notification_template.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()
	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	notificationTemplate = string(b)
}

func SetConfig(token string, url string) {
	config = Config{NotificationsToken: token, NotificationsURL: url}
}

func SendNotifications(data SendData) error {
	tmpl, err := template.New("").Parse(notificationTemplate)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return err
	}
	message := buffer.String()

	var wg sync.WaitGroup
	wg.Add(len(data.ChatID))
	for _, chatID := range data.ChatID {
		chatID := chatID
		go func() {
			defer wg.Done()

			chatIntID, err := strconv.ParseInt(chatID, 10, 64)
			if err != nil {
				log.Println(fmt.Sprintf("bad chat_id: %s", err))
				return
			}

			err = sendNotification(chatIntID, message)
			if err != nil {
				log.Println(fmt.Sprintf("unable to send notification to %s: %s", chatID, err.Error()))
			}
		}()
	}
	wg.Wait()

	return nil
}

func sendNotification(chatID int64, message string) error {
	jsonData, err := json.Marshal(request{
		ChatID:  chatID,
		Message: message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.NotificationsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.NotificationsToken)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if string(resBody) != "OK" {
		return errors.New(string(resBody))
	}

	return nil
}
