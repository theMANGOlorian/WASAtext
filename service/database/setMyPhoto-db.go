package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
)

func (db *appdbimpl) SetMyPhoto(id string) (string, error) {

	const query1 = `UPDATE users SET photo = ? WHERE id = ?`
	const query2 = `SELECT 1 FROM users WHERE id = ?`
	const query3 = `
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
	err := db.c.QueryRow(query2, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not exists")
		}
		return "", err
	}

	tx, err := db.GetTx()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %v", err)
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
		err = tx.QueryRow(query3, codeImg, codeImg, codeImg).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.Exec(query1, codeImg, id)
			if err != nil {
				return "", err
			}
			err := db.CloseTx(tx, true)
			if err != nil {
				return "", err
			}
			break
		}
		if err != nil {
			return "", err
		}
		codeImg = ksuid.New().String()
	}
	return codeImg, nil
}
