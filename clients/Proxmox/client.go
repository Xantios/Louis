package Proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	SimpleLogger "git.sacredheart.it/xantios/simple-logger"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client   *http.Client
	url      string
	username string
	password string
	nodeName string
	logger   *SimpleLogger.SimpleLogger
}

type LoginResponse struct {
	Data struct {
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Ticket              string `json:"ticket"`
		Username            string `json:"username"`
	} `json:"data"`
}

type UpdateResponse struct {
	Data []struct {
		Package     string `json:"package"`
		Version     string `json:"version"`
		OldVersion  string `json:"old-version"`
		Description string `json:"description"`
	} `json:"data"`
}

// Feel free to have an opinion on this, send patches etc
// @TODO: Make configurable
var insecureClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func New(logger *SimpleLogger.SimpleLogger, pveName string, url string, username string, password string) *Client {
	return &Client{
		client:   insecureClient,
		url:      url,
		username: username,
		password: password,
		nodeName: pveName,
		logger:   logger,
	}
}

func login(logger *SimpleLogger.SimpleLogger, host string, username, password string) (ticket string, err error) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	resp, err := insecureClient.PostForm(host+"/access/ticket", form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// fmt.Printf("Expected 200 got: %s\n", resp.Status)
		logger.Error("Expected 200 got: %s", resp.Status)

		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		resp.Body.Close()
		// fmt.Printf("Body: %s\n", string(body))

		return "", fmt.Errorf("status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(data, &loginResp); err != nil {
		return "", err
	}

	return loginResp.Data.Ticket, nil
}

func checkUpdates(host string, nodeName string, ticket string) (int, error) {

	req, err := http.NewRequest("GET", host+"/nodes/"+nodeName+"/apt/update", nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Cookie", "PVEAuthCookie="+ticket)
	req.Header.Set("Accept", "application/json")

	resp, err := insecureClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("status: %s, body: %s", resp.Status, string(body))
	}

	var updateResp UpdateResponse
	if err = json.Unmarshal(body, &updateResp); err != nil {
		return 0, err
	}

	if len(updateResp.Data) == 0 {
		return 0, nil
	}

	return len(updateResp.Data), nil
}

func (c *Client) Update() (bool, string, error) {
	ticket, err := login(c.logger, c.url, c.username, c.password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	c.logger.Debug("Ticket collected from PVE API")

	count, err := checkUpdates(c.url, c.nodeName, ticket)
	if err != nil {
		c.logger.Error("Failed to check for updates: %v", err)
	}

	updateAvailable := count > 0
	return updateAvailable, fmt.Sprintf("Found %d updates for %s", count, c.nodeName), nil
}

func (c *Client) Backup() error {
	return nil
}
