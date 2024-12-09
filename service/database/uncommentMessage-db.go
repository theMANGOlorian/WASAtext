package database

import (
	"fmt"
	"strings"
)

func (db *appdbimpl) UncommentMessage(userId string, messageId string) (int, error) {

	const query1 = `DELETE FROM reactions WHERE owner = ? AND messageId = ?`
	const query2 = `SELECT EXISTS(SELECT 1 FROM reactions WHERE owner = ? AND messageId = ?)`

	var exists int
	err := db.c.QueryRow(query2, userId, messageId).Scan(&exists)
	if err != nil {
		return 500, err
	}
	if exists == 0 {
		return 404, fmt.Errorf("not found")
	}

	_, err = db.c.Exec(query1, userId, messageId)
	if err != nil {
		if strings.HasPrefix(err.Error(), "NOT FOUND") {
			return 404, err
		}
		return 500, err
	}

	return 204, err
}
