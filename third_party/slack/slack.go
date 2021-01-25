package slack

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/midnight-trigger/todo/third_party/http_client"

	"github.com/spf13/viper"
)

type Payload struct {
	IconEmoji   string        `json:"icon_emoji"`
	Channel     string        `json:"channel"`
	Username    string        `json:"username"`
	Attachments []*Attachment `json:"attachments"`
}

type Attachment struct {
	Color    string   `json:"color"`
	Pretext  string   `json:"pretext"`
	Text     string   `json:"text"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

// send method : Create JSON response data
func send(URL string, Method string, Code int, header map[string][]string, requestPayload string, errorType string, message string, errStack string, rawError string) (err error) {

	host, _ := os.Hostname()

	payload := new(Payload)
	payload.IconEmoji = ":scream_cat:"
	payload.Username = host
	payload.Channel = viper.GetString("slack.channel")

	if len(payload.Channel) > 0 {
		attachment := new(Attachment)
		pretext := []string{"<!here>", "[URL] " + Method + ": " + URL}
		attachment.Pretext = strings.Join(pretext, "\n")
		payload.Attachments = append(payload.Attachments, attachment)

		attachment = new(Attachment)
		attachment.Color = "danger"
		attachment.MrkdwnIn = []string{"text", "pretext"}
		attachment.Pretext = "[" + errorType + "]"
		attachment.Text = message
		payload.Attachments = append(payload.Attachments, attachment)

		attachment = new(Attachment)
		attachment.Color = "danger"
		attachment.MrkdwnIn = []string{"text", "pretext"}
		attachment.Pretext = "[raw error message]"
		attachment.Text = rawError
		payload.Attachments = append(payload.Attachments, attachment)

		attachment = new(Attachment)
		attachment.Color = "danger"
		attachment.MrkdwnIn = []string{"text", "pretext"}
		attachment.Pretext = "[Stack Trace]"
		attachment.Text = errStack
		payload.Attachments = append(payload.Attachments, attachment)

		attachment = new(Attachment)
		attachment.Color = "#337ab7"
		attachment.MrkdwnIn = []string{"text", "pretext"}
		attachment.Pretext = "[Request Payload]"
		attachment.Text = requestPayload
		payload.Attachments = append(payload.Attachments, attachment)

		attachment = new(Attachment)
		attachment.Color = "#337ab7"
		attachment.MrkdwnIn = []string{"text", "pretext"}
		attachment.Pretext = "[Request Headers]"
		b := new(bytes.Buffer)
		for key, value := range header {
			fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
		}
		attachment.Text = b.String()
		payload.Attachments = append(payload.Attachments, attachment)

		client := new(http_client.HTTPClient)
		client.Body = payload
		client.URI = viper.GetString("slack.endpoint")
		//client.UseProxy = true
		_, err = client.Post()
	}
	return
}

// SlackSend method : Create JSON response data
func SlackSend(URL string, Method string, Code int, header map[string][]string, requestPayload string, errorType string, message string, errStack string, rawError string) {
	go send(URL, Method, Code, header, requestPayload, errorType, message, errStack, rawError)
}
