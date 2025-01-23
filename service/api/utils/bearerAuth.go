package utils

import (
	"fmt"
	"log"
	"strings"
)

func CheckAuthorizationField(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "

	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	if !strings.Contains(authHeader, bearerPrefix) {
		log.Println("AuthHeader: ", authHeader, "\nBearerPrefix:", bearerPrefix)
		return "", fmt.Errorf("invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if !CheckIdentifier(token) {
		return "", fmt.Errorf("invalid token")
	}

	return token, nil
}
