package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) LeaveGroup(userId string, groupId string) (int, error) {

	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ?`
	const query4 = `DELETE FROM members WHERE userId = ? AND conversationId = ?`

	var exists int
	err := db.c.QueryRow(query1, userId).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return 404, err
	}

	err = db.c.QueryRow(query2, groupId).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return 404, err
	}

	err = db.c.QueryRow(query3, userId, groupId).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return 404, err
	}

	_, err = db.c.Exec(query4, userId, groupId)
	if err != nil {
		return 500, err
	}

	return 200, nil

}
