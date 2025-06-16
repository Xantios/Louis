package Webhook

import (
	"bytes"
	"encoding/json"
	SimpleLogger "git.sacredheart.it/xantios/simple-logger"
	"net/http"
)

type Webhook struct {
	hookUrl string
	logger  *SimpleLogger.SimpleLogger
}

func New(logger *SimpleLogger.SimpleLogger, hookUrl string) *Webhook {
	return &Webhook{
		hookUrl: hookUrl,
		logger:  logger,
	}
}

func (s *Webhook) Send(message string) error {

	m, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Failed to marshal JSON: %v", err)
		return err
	}

	resp, err := http.Post(s.hookUrl, "application/json", bytes.NewBuffer(m))
	if err != nil {
		s.logger.Error("Failed to send request: %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
