package database

import (
	"WASAtext/service/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"strings"
)

func (db *appdbimpl) SendImage(userId string, conversationId string) (int, *utils.SendMessageResponseBody, error) {

	const query1 = `SELECT username FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query4 = `INSERT INTO messages(id, type, photo, conversation, status) VALUES (?,'photo',?,?,'none')`
	const query5 = `SELECT timestamp FROM messages WHERE id = ? LIMIT 1`
	const query6 = `
					SELECT 1
					FROM conversations c
					WHERE c.photo = ?
					UNION
					SELECT 1
					FROM users u
					WHERE u.photo = ?
					UNION
					SELECT 1
					FROM messages m
					WHERE m.photo = ?
					LIMIT 1;
					`

	var exists int
	var username string
	err := db.c.QueryRow(query1, userId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, nil, fmt.Errorf("user not exists")
		}
		return 500, nil, err
	}

	err = db.c.QueryRow(query2, conversationId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, nil, fmt.Errorf("conversation not exists")
		}
		return 500, nil, err
	}

	err = db.c.QueryRow(query3, userId, conversationId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 403, nil, fmt.Errorf("user cannot send message in this conversation because he is not a member")
		}
		return 500, nil, err
	}

	tx, err := db.GetTx()
	if err != nil || tx == nil {
		return 500, nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			// Se un errore è presente, esegui il rollback
			rollbackErr := db.CloseTx(tx, false)
			if rollbackErr != nil {
				// In caso di errore durante il rollback, registriamo un errore
				fmt.Println("Failed to rollback transaction:", rollbackErr)
			}
		}
	}()

	var messageId string
	var codeImg string
	var violation bool
	for {
		messageId = uuid.New().String()
		codeImg = ksuid.New().String()
		violation = false

		err = tx.QueryRow(query6, codeImg, codeImg, codeImg).Scan(&exists)
		if err == nil && exists == 1 {
			violation = true
			continue
		}
		if err != nil && !(errors.Is(err, sql.ErrNoRows)) {
			return 500, nil, err
		}

		_, err = tx.Exec(query4, messageId, codeImg, conversationId)
		if err != nil && !(strings.Contains(err.Error(), "UNIQUE constraint failed")) {
			return 500, nil, err
		}
		if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			violation = true
		}
		if violation == false {
			err = nil
			err := db.CloseTx(tx, true)
			if err != nil {
				return 500, nil, err
			}
			break
		}

	}
	var ts string
	err = db.c.QueryRow(query5, messageId).Scan(&ts)
	if err != nil {
		return 500, nil, err
	}

	var response = &utils.SendMessageResponseBody{
		MessageId:   messageId,
		SenderId:    userId,
		Username:    username,
		Text:        "",
		Image:       codeImg,
		ReplyTo:     "",
		Timestamp:   ts,
		Status:      "none",
		TypeContent: "photo",
		Reactions:   []utils.Reactions{},
	}

	return 201, response, nil
}
