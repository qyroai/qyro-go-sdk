package qyro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type QyroClient struct {
	BaseURL string
	Token   string
	Timeout time.Duration
}

func NewQyroClient(baseURL, token string, timeout time.Duration) (*QyroClient, error) {
	if baseURL == "" {
		return nil, &ConfigurationError{"base_url is required"}
	}
	return &QyroClient{
		BaseURL: baseURL,
		Token:   token,
		Timeout: timeout,
	}, nil
}

func (c *QyroClient) url(path string) string {
	return fmt.Sprintf("%s%s", c.BaseURL, path)
}

func (c *QyroClient) clientHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
}

func (c *QyroClient) CreateSession(assistantID string, context map[string]interface{}) (*Session, error) {
	path := fmt.Sprintf("/client/api/v1/assistants/%s/sessions", assistantID)
	payload, _ := json.Marshal(map[string]interface{}{"context": context})
	req, _ := http.NewRequest("POST", c.url(path), bytes.NewBuffer(payload))
	c.clientHeaders(req)

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

func (c *QyroClient) FetchSessionMessages(assistantID, sessionID string) ([]Message, error) {
	path := fmt.Sprintf("/client/api/v1/assistants/%s/sessions/%s/messages", assistantID, sessionID)
	req, _ := http.NewRequest("GET", c.url(path), nil)
	c.clientHeaders(req)

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

func (c *QyroClient) Chat(assistantID, sessionID, message string) ([]Message, error) {
	path := fmt.Sprintf("/client/api/v1/assistants/%s/sessions/%s/chat", assistantID, sessionID)
	payload, _ := json.Marshal(map[string]string{"message": message})
	req, _ := http.NewRequest("POST", c.url(path), bytes.NewBuffer(payload))
	c.clientHeaders(req)

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