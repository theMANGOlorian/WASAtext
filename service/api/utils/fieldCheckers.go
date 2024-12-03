package utils

import (
	"github.com/google/uuid"
)

func CheckIdentifier(identifier string) bool {
	_, err := uuid.Parse(identifier)
	if err != nil {
		return false
	}
	return true
}
