package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dark-warlord14/CVENotifier/internal/errors"
	"github.com/dark-warlord14/CVENotifier/internal/util"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func NotifySlack(vulnTitle string, link string, published string, categories string, description string, slackWebhook string) error {
	description = util.RemoveHTMLTags(description)

	message := SlackMessage{
		Text: "Title: " + vulnTitle + "\nLink: " + link + "\nDate Published: " + published + "\nDescription: " + description,
	}

	// Encode message payload as JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return &errors.SlackNotificationError{Message: "Failed to marshal message: " + err.Error()}
	}

	// Make POST request to Slack webhook
	resp, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return &errors.SlackNotificationError{Message: "Failed to send message, check if slack webhook is valid: " + err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &errors.SlackNotificationError{Message: fmt.Sprintf("Failed to send message, status code: %d", resp.StatusCode)}
	}

	return nil
}
