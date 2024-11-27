package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
)

// Get userId, if user not exists then it will create
func (db *appdbimpl) DoLogin(username string) (string, error) {
	var id string

	// Check if user exists
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err == nil {
		// Utente esiste già, restituisci l'ID trovato
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("error while checking username: %w", err)
	}

	// User not found, it will created
	for {
		// Generate an UUID
		newId := uuid.New().String()

		// Trying to insert the new user into database
		_, err = db.c.Exec("INSERT INTO users (id, username) VALUES (?, ?)", newId, username)
		if err == nil {
			// Inserimento riuscito, restituisci l'ID generato
			return newId, nil
		}

		return "", fmt.Errorf("error while creating a new user: %w", err)
	}
}

// update username in the database using userId
func (db *appdbimpl) SetMyUserName(id string, newUsername string) (string, error) {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE  id = ?", newUsername, id)
	if err != nil {
		return "", fmt.Errorf("error while updating username: %w", err)
	}
	return newUsername, nil
}

func (db *appdbimpl) SetMyPhoto(id string) (string, error) {

	codeImg := ksuid.New().String()
	_, err := db.c.Exec("UPDATE users SET photo = ? WHERE id = ?", codeImg, id)
	if err != nil {
		return "", fmt.Errorf("error while updating photo profile: %w", err)
	}
	return codeImg, nil
}
