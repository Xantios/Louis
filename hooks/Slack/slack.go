package Slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

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
