package OPNSense

import (
	SimpleLogger "git.sacredheart.it/xantios/simple-logger"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Client     *http.Client
	backupPath string
	base       string
	username   string
	password   string
	debug      bool
	logger     *SimpleLogger.SimpleLogger
}

func New(logger *SimpleLogger.SimpleLogger, baseURL string, backupPath string, username string, password string, debug bool) *Client {

	c := &http.Client{
		// Timeout is high because of VPN
		Timeout: 30 * time.Second, // @TODO: Make this configurable
	}

	return &Client{
		Client:     c,
		base:       baseURL,
		username:   username,
		password:   password,
		backupPath: backupPath,
		debug:      debug,
		logger:     logger,
	}
}

func (c *Client) Get(url string) (*http.Response, string, error) {

	c.logger.Debugf("Sending REQ to :: %s", c.base+url)

	req, err := http.NewRequest("GET", c.base+url, nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(c.username, c.password)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	resp.Body.Close()

	return resp, string(body), nil
}

func (c *Client) Post(url string) (*http.Response, string, error) {
	req, err := http.NewRequest("POST", c.base+url, nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(c.username, c.password)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	resp.Body.Close()

	return resp, string(body), nil
}
