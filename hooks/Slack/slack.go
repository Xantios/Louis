package Slack

import (
	"bytes"
	"encoding/json"
	SimpleLogger "git.sacredheart.it/xantios/simple-logger"
	"net/http"
)

// SlackPayload defines the structure of the message payload
type SlackPayload struct {
	Text string `json:"text"`
}

type Slack struct {
	hookUrl string
	logger  *SimpleLogger.SimpleLogger
}

func New(logger *SimpleLogger.SimpleLogger, hookUrl string) *Slack {
	return &Slack{
		hookUrl: hookUrl,
		logger:  logger,
	}
}

func (s *Slack) Send(message string) error {

	msg := SlackPayload{
		Text: message,
	}

	m, err := json.Marshal(msg)
	if err != nil {
		s.logger.Error("Failed to marshal message to JSON: %v", err)
		return err
	}

	resp, err := http.Post(s.hookUrl, "application/json", bytes.NewBuffer(m))
	if err != nil {
		s.logger.Error("Failed to send request to Slack: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
