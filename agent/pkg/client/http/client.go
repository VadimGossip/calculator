package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client interface {
	SendGetRequest(endpoint string, request, response interface{}) error
	SendPostRequest(endpoint string, request, response interface{}) error
}

type client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL string, apiKey string, ttl time.Duration) *client {
	return &client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: ttl,
		},
	}
}

func (c *client) SendGetRequest(endpoint string, request, response interface{}) error {
	url := c.baseURL + endpoint

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Close = true

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) SendPostRequest(endpoint string, request, response interface{}) error {
	url := c.baseURL + endpoint

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Close = true

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}
