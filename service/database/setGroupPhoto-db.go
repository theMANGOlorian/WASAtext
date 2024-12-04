package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"strings"
)

func (db *appdbimpl) SetGroupPhoto(userId string, groupId string) (string, int, error) {

	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? AND type = 'group' LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query4 = `UPDATE conversations SET photo = ? WHERE id = ?`
	const query5 = `SELECT 1 FROM conversations WHERE photo = ?`

	var exists int
	err := db.c.QueryRow(query1, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 404, fmt.Errorf("user not exists")
		}
		return "", 500, err
	}

	err = db.c.QueryRow(query2, groupId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 404, fmt.Errorf("group not exists")
		}
		return "", 500, err
	}

	err = db.c.QueryRow(query3, userId, groupId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 403, fmt.Errorf("user is not a member of group")
		}
		return "", 500, err
	}

	codeImg := ksuid.New().String()
	for {
		_, err := db.c.Exec(query4, codeImg, groupId)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			codeImg = ksuid.New().String()
		} else {
			return "", 500, err
		}

	}

	return codeImg, 201, nil

}
