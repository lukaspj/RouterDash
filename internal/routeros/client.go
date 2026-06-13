package routeros

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	address    string
	username   string
	password   string
	httpClient *http.Client
}

func NewClient(address, username, password string, insecureSkipVerify bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	return &Client{
		address:  strings.TrimRight(address, "/"),
		username: username,
		password: password,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

func (c *Client) Get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/%s", c.address, strings.TrimLeft(path, "/"))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("routeros api error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
