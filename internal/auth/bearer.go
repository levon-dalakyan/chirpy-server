package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with Bearer")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}
