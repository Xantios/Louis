package Webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Webhook struct {
	hookUrl string
}

func New(hookUrl string) *Webhook {
	return &Webhook{
		hookUrl: hookUrl,
	}
}

func (s *Webhook) Send(message string) error {

	m, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		return err
	}

	resp, err := http.Post(s.hookUrl, "application/json", bytes.NewBuffer(m))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
