package database

import (
	"github.com/segmentio/ksuid"
	"strings"
)

func (db *appdbimpl) SetMyPhoto(id string) (string, error) {

	const query1 = `UPDATE users SET photo = ? WHERE id = ?`
	codeImg := ksuid.New().String()
	for {
		_, err := db.c.Exec(query1, codeImg, id)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			codeImg = ksuid.New().String()
		} else {
			return "", err
		}

	}
	return codeImg, nil
}
