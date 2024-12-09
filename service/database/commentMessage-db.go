package database

import (
	"github.com/google/uuid"
	"strings"
)

func (db *appdbimpl) CommentMessage(userId string, messageId string, reaction string) (int, error) {

	// c'è un trigger su reactions INSERT
	const query1 = `INSERT INTO reactions (id, owner, messageId, emoji) VALUES (?,?,?,?)`
	const query2 = `SELECT EXISTS (SELECT 1 FROM reactions WHERE owner = ? and messageId = ?)`
	const query3 = `UPDATE reactions SET emoji = ? WHERE owner = ? and messageId = ?`
	var exists int
	err := db.c.QueryRow(query2, userId, messageId).Scan(&exists)
	if err != nil {
		return 500, err
	}
	if exists == 0 {
		reactionId := uuid.New().String()
		for {
			_, err = db.c.Exec(query1, reactionId, userId, messageId, reaction)
			if err != nil {
				if strings.HasPrefix(err.Error(), "UNIQUE") {
					reactionId = uuid.New().String()
					continue
				}

				return 500, err
			}
			break
		}
	} else {
		_, err = db.c.Exec(query3, reaction, userId, messageId)
		if err != nil {
			return 500, err
		}
	}
	return 204, nil

}
