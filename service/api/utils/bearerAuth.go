package utils

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func CheckAuthorizationField(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "

	if authHeader == "" {
		return "", fmt.Errorf("Missing Authorization header")
	}
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("Invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	_, err := uuid.Parse(token)
	if err != nil {
		return "", fmt.Errorf("Invalid Token")
	}

	return token, nil
}
