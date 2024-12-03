package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) AddToGroupPermission(userId string, groupId string) (int, error) {
	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? AND type = 'group' LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`

	var exists int
	err := db.c.QueryRow(query1, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("user not exists")
		}
		return 500, err
	}

	err = db.c.QueryRow(query2, groupId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("group not exists")
		}
		return 500, err
	}

	err = db.c.QueryRow(query3, userId, groupId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 403, fmt.Errorf("user is not a member of group")
		}
		return 500, err
	}

	return 200, nil
}

func (db *appdbimpl) AddToGroup(friendId string, groupId string) (int, error) {
	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query4 = `INSERT INTO members (userId,conversationId) VALUES (?,?)`

	var exists int
	err := db.c.QueryRow(query1, friendId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, fmt.Errorf("friend not exists")
		}
		return 500, err
	}
	err = db.c.QueryRow(query3, friendId, groupId).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.c.Exec(query4, friendId, groupId)
		if err != nil {
			return 500, err
		}
		return 201, nil
	} else if err == nil {
		return 409, fmt.Errorf("friend already a member")
	} else {
		return 500, err
	}

}
