package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with ApiKey")
	}

	apiKey := strings.TrimPrefix(authHeader, prefix)
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return "", errors.New("ApiKey is empty")
	}

	return apiKey, nil
}
