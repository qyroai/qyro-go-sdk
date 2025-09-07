package qyro

import (
	"fmt"
	"net/http"
)

type QyroError struct {
	Message string
}

func (e *QyroError) Error() string {
	return e.Message
}

type HTTPError struct {
	StatusCode int
	Body       string
	Response   *http.Response
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

type ConfigurationError struct {
	Msg string
}

func (e *ConfigurationError) Error() string {
	return e.Msg
}
