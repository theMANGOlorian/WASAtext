package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
)

func (db *appdbimpl) SetMyPhoto(id string) (string, error) {

	const query1 = `UPDATE users SET photo = ? WHERE id = ?`
	codeImg := ksuid.New().String()
	for {
		_, err := db.c.Exec(query1, codeImg, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			return "", fmt.Errorf("error while updating photo profile: %w", err)
		}
		codeImg = ksuid.New().String()
	}
	return codeImg, nil
}
