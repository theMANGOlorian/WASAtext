package utils

import (
	"regexp"
)

func CheckIdentifier(identifier string) bool {
	const pattern = "^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	re := regexp.MustCompile(pattern)
	return re.MatchString(identifier)

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
	return re.MatchString(photoId)
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
