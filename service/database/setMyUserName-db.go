package database

import "fmt"

// update username in the database using userId
func (db *appdbimpl) SetMyUserName(id string, newUsername string) (string, error) {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE  id = ?", newUsername, id)
	if err != nil {
		return "", fmt.Errorf("error while updating username: %w", err)
	}
	return newUsername, nil
}
