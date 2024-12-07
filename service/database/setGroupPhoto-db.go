package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
)

func (db *appdbimpl) SetGroupPhoto(userId string, groupId string) (string, int, error) {

	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? AND type = 'group' LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query4 = `UPDATE conversations SET photo = ? WHERE id = ?`
	const query5 = `
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

	tx, err := db.GetTx()
	if err != nil {
		return "", 500, fmt.Errorf("failed to begin transaction: %v", err)
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

	codeImg := ksuid.New().String()
	for {
		err = tx.QueryRow(query5, codeImg, codeImg, codeImg).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.Exec(query4, codeImg, groupId)
			if err != nil {
				return "", 500, err
			}
			err := db.CloseTx(tx, true)
			if err != nil {
				return "", 500, err
			}
			break
		}
		if err != nil {
			return "", 500, err
		}
		codeImg = ksuid.New().String()
	}
	return codeImg, 201, nil

}
