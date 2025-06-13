package Slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// curl -X POST --data-urlencode "payload={\"channel\": \"#opnsense\", \"username\": \"webhookbot\", \"text\": \"This is posted to #my-channel-here and comes from a bot named webhookbot.\", \"icon_emoji\": \":ghost:\"}" https://hooks.slack.com/services/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

// SlackPayload defines the structure of the message payload
type SlackPayload struct {
	Text string `json:"text"`
}

type Slack struct {
	hookUrl string
}

func New(hookUrl string) *Slack {
	return &Slack{
		hookUrl: hookUrl,
	}
}

func (s *Slack) Send(message string) error {

	msg := SlackPayload{
		Text: message,
	}

	m, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		return err
	}

	resp, err := http.Post(s.hookUrl, "application/json", bytes.NewBuffer(m))
	if err != nil {
		log.Fatalf("Failed to send request to Slack: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
