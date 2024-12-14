package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func (db *appdbimpl) DoLogin(username string) (string, string, error) {
	var id string
	var photo sql.NullString

	// Check if user exists
	err := db.c.QueryRow("SELECT id, photo FROM users WHERE username = ?", username).Scan(&id, &photo)
	if err == nil {
		if photo.Valid {
			photoCode := photo.String
			return id, photoCode, nil
		} else {
			return id, "", nil
		}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "","", fmt.Errorf("error while checking username: %w", err)
	}

	// User not found, it will created
	for {
		// Generate an UUID
		newId := uuid.New().String()

		// Trying to insert the new user into database
		_, err = db.c.Exec("INSERT INTO users (id, username) VALUES (?, ?)", newId, username)
		if err == nil {
			return newId, "", nil
		}
	}




}
