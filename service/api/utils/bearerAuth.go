package utils

import (
	"fmt"
	"strings"
)

func CheckAuthorizationField(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "

	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if !CheckIdentifier(token) {
		return "", fmt.Errorf("invalid Token")
	}

	return token, nil
}
