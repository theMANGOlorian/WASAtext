package database

import (
	"fmt"
	"github.com/google/uuid"
	"database/sql"
)

func (db *appdbimpl) DoLogin(name string) (string, error) {
	var id string

	// Controlla se lo username esiste già
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", name).Scan(&id)
	if err == nil {
		// Utente esiste già, restituisci l'ID trovato
		return id, nil
	}
	if err != nil && err != sql.ErrNoRows {
		// Errore imprevisto durante la query
		return "", fmt.Errorf("errore durante il controllo dello username: %w", err)
	}

	// Utente non trovato, creane uno nuovo
	for {
		// Genera un nuovo UUID
		newId := uuid.New().String()

		// Prova a inserire il nuovo utente nel database
		_, err = db.c.Exec("INSERT INTO users (id, username) VALUES (?, ?)", newId, name)
		if err == nil {
			// Inserimento riuscito, restituisci l'ID generato
			return newId, nil
		}
		// Altro errore: restituisci l'errore
		return "", fmt.Errorf("errore durante la creazione dell'utente: %w", err)
	}
}
