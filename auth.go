package qyro

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ApiKeyAuth struct {
	ApiKeyID     string
	ApiKeySecret string
}

func NewApiKeyAuth(apiKeyID, apiKeySecret string) *ApiKeyAuth {
	return &ApiKeyAuth{
		ApiKeyID:     apiKeyID,
		ApiKeySecret: apiKeySecret,
	}
}

func (a *ApiKeyAuth) HeaderValue() string {
	return "ApiKey " + a.ApiKeySecret
}

type ClientTokenGenerator struct {
	ApiKeyID     string
	ApiKeySecret string
}

func NewClientTokenGenerator(apiKeyID, apiKeySecret string) *ClientTokenGenerator {
	return &ClientTokenGenerator{
		ApiKeyID:     apiKeyID,
		ApiKeySecret: apiKeySecret,
	}
}

func (c *ClientTokenGenerator) Generate(context map[string]interface{}) (string, error) {
	subjectBytes, err := json.Marshal(context)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"sub":  string(subjectBytes),
		"iat":  now,
		"exp":  now + 24*30*3600,
		"type": "client",
		"iss":  c.ApiKeyID,
		"aud":  "qyro",
		"jti":  uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = c.ApiKeyID

	return token.SignedString([]byte(c.ApiKeySecret))
}