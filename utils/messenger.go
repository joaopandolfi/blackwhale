package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/remotes/request"
)

type DbgMessage struct {
	Level    string      `json:"level"`
	Datetime string      `json:"datetime"`
	Service  string      `json:"service"`
	Message  string      `json:"message"`
	Context  interface{} `json:"context"`
}

// Payload to send Attachment on slack
type SlackAttachment struct {
	Text        string              `json:"text"`
	Channel     string              `json:"channel"`
	Token       string              `json:"token"`
	Username    string              `json:"username"`
	Attachments []map[string]string `json:"attachments"`
}

var count int

func newBaseMessage(level string, message string, data interface{}) DbgMessage {
	return DbgMessage{
		Level:    level,
		Datetime: time.Now().UTC().Format(time.RFC3339),
		Service:  configurations.Configuration.Name,
		Message:  message,
		Context:  data,
	}
}

func dispatch(m DbgMessage) {
	r, _ := json.Marshal(m)
	fmt.Println(string(r))
}

func levelToColor(level string) (color string) {
	switch level {
	case "INFO":
		color = "#FFFF00"
		break

	case "DEBUG":
		color = "#36a64f"
		break

	case "ERROR":
		color = "#DF3A01"
		break
	}
	return
}

func levelToEmogi(level string) (emogi string) {
	switch level {
	case "INFO":
		emogi = ":information_source:"
		break

	case "DEBUG":
		emogi = ":white_check_mark:"
		break

	case "ERROR":
		emogi = ":exclamation:"
		break
	}
	return
}

func slackDispatch(m DbgMessage) {
	var url string
	defer func() {
		if r := recover(); r != nil {
			me, _ := json.Marshal(newBaseMessage("ERROR", "[Messager][SlackDispatch] - Error on slack dispatch URL:"+url, r))
			message := string(me)
			fmt.Println(message)
		}
	}()

	content, _ := json.Marshal(m.Context)

	var attachment []map[string]string
	attachment = append(attachment, map[string]string{
		"title": fmt.Sprintf("%s %s", levelToEmogi(m.Level), m.Message),
		"text":  fmt.Sprintf("*Timestamp:* %s \n*Context:* %s", m.Datetime, string(content)),
		"color": levelToColor(m.Level),
	})

	payload2 := SlackAttachment{
		Channel:     configurations.Configuration.SlackChannel,
		Token:       configurations.Configuration.SlackToken,
		Username:    "blackwhale",
		Attachments: attachment,
	}

	body, _ := json.Marshal(payload2)

	url = getSlackUrl()
	request.Post(url, body)

}

func getSlackUrl() string {
	if count > 3 {
		count = 0
	} else {
		count++
	}

	return configurations.Configuration.SlackWebHook[count]
}

func Info(message string, data ...interface{}) {
	dbg := newBaseMessage("INFO", message, data)
	dispatch(dbg)
	//go slackDispatch(dbg)
}

func Feedback(message string, data ...interface{}) {
	dbg := newBaseMessage("DEBUG", message, data)
	dispatch(dbg)
	//go slackDispatch(dbg)
}

func Debug(message string, data ...interface{}) {
	dbg := newBaseMessage("DEBUG", message, data)
	dispatch(dbg)
}

func Error(message string, data ...interface{}) {
	dbg := newBaseMessage("ERROR", message, data)
	dispatch(dbg)
}

func CriticalError(message string, data ...interface{}) {
	dbg := newBaseMessage("ERROR", message, data)
	dispatch(dbg)
	//go slackDispatch(dbg)
}

func Logger(data ...interface{}) error {
	dbg := newBaseMessage("DEBUG", "LOGGER", data)
	dispatch(dbg)
	return nil
}
