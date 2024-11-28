package database

import (
	"fmt"
	"github.com/segmentio/ksuid"
)

func (db *appdbimpl) SetMyPhoto(id string) (string, error) {

	codeImg := ksuid.New().String()
	_, err := db.c.Exec("UPDATE users SET photo = ? WHERE id = ?", codeImg, id)
	if err != nil {
		return "", fmt.Errorf("error while updating photo profile: %w", err)
	}
	return codeImg, nil
}
