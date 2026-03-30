package auth

import (
	"errors"
	"net/http"
	"strings"
)

// get token from authorization header
func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("Invalid header")
	}
	// check if header starts with bearer
	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("Invalid header")
	}
	// strip "Bearer "
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" {
		return "", errors.New("Invalid token")
	}
	return token, nil
}
