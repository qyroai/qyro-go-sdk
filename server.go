package qyro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type QyroServerClient struct {
	BaseURL string
	Auth    *ApiKeyAuth
	Timeout time.Duration
}

func NewQyroServerClient(baseURL, apiKeyID, apiKeySecret string, timeout time.Duration) (*QyroServerClient, error) {
	if baseURL == "" {
		return nil, &ConfigurationError{"base_url is required"}
	}
	return &QyroServerClient{
		BaseURL: baseURL,
		Auth:    NewApiKeyAuth(apiKeyID, apiKeySecret),
		Timeout: timeout,
	}, nil
}

func (c *QyroServerClient) url(path string) string {
	return fmt.Sprintf("%s%s", c.BaseURL, path)
}

func raiseForStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(body), Response: resp}
	}
	return nil
}

func (c *QyroServerClient) CreateSession(assistantID string, context map[string]interface{}) (*Session, error) {
	path := fmt.Sprintf("/server/api/v1/assistants/%s/sessions", assistantID)
	payload, _ := json.Marshal(map[string]interface{}{"context": context})
	req, _ := http.NewRequest("POST", c.url(path), bytes.NewBuffer(payload))
	req.Header.Set("Authorization", c.Auth.HeaderValue())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := raiseForStatus(resp); err != nil {
		return nil, err
	}

	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (c *QyroServerClient) FetchSessionMessages(assistantID, sessionID string) ([]Message, error) {
	path := fmt.Sprintf("/server/api/v1/assistants/%s/sessions/%s/messages", assistantID, sessionID)
	req, _ := http.NewRequest("GET", c.url(path), nil)
	req.Header.Set("Authorization", c.Auth.HeaderValue())

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := raiseForStatus(resp); err != nil {
		return nil, err
	}

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *QyroServerClient) Chat(assistantID, sessionID, message string) ([]Message, error) {
	path := fmt.Sprintf("/server/api/v1/assistants/%s/sessions/%s/chat", assistantID, sessionID)
	payload, _ := json.Marshal(map[string]string{"message": message})
	req, _ := http.NewRequest("POST", c.url(path), bytes.NewBuffer(payload))
	req.Header.Set("Authorization", c.Auth.HeaderValue())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := raiseForStatus(resp); err != nil {
		return nil, err
	}

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}