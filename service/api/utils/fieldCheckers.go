package utils

import (
	"github.com/google/uuid"
	"regexp"
)

func CheckIdentifier(identifier string) bool {
	_, err := uuid.Parse(identifier)
	if err != nil {
		return false
	}
	return true
}

func CheckName(name string) bool {
	const pattern = "^[0-9a-zA-Z]*$"
	re := regexp.MustCompile(pattern)
	if len(name) >= 3 && len(name) <= 25 && re.MatchString(name) {
		return true
	}
	return false
}

func CheckPhotoId(photoId string) bool {
	const pattern = "^[0-9a-zA-Z]*$"
	re := regexp.MustCompile(pattern)
	if re.MatchString(photoId) {
		return true
	}
	return false

}

func CheckReactions(reaction string) bool {
	reactionsAllowed := []string{"😊", "😂", "😍", "😎", "😭"}
	for _, r := range reactionsAllowed {
		if reaction == r {
			return true
		}
	}
	return false
}
