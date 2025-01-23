package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func (db *appdbimpl) StartConversationPrivate(id string, friendName string) (string, error) {
	// Inizia la transazione
	tx, err := db.GetTx()
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			// Se un errore è presente, esegui il rollback
			rollbackErr := db.CloseTx(tx, false)
			if rollbackErr != nil {
				// In caso di errore durante il rollback, registriamo un errore
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
		}
	}()

	// Recupera l'ID dell'amico
	var friendId string
	err = tx.QueryRow("SELECT id FROM users WHERE username = ?", friendName).Scan(&friendId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("NOT FOUND")
		}
		return "", fmt.Errorf("error while checking friend name: %w", err)
	}

	// Controlla se esiste già una conversazione privata
	var count int
	query := "SELECT COUNT(*) FROM members m1 JOIN members m2 ON m1.conversationId = m2.conversationId JOIN conversations c ON m1.conversationId = c.id WHERE c.type = 'private' AND m1.userId = ? AND m2.userId = ?;"
	err = tx.QueryRow(query, id, friendId).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("error checking if users are in the same conversation: %w", err)
	}
	if count == 1 {
		return "", fmt.Errorf("EXISTS")
	}

	// Genera un nuovo ID per la conversazione
	var conversationId string
	for {
		conversationId = uuid.New().String()
		err = tx.QueryRow("SELECT id FROM conversations WHERE id = ?", conversationId).Scan(&conversationId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err = tx.Exec("INSERT INTO conversations (id, type) VALUES (?, 'private')", conversationId)
				if err != nil {
					return "", fmt.Errorf("error while inserting conversation: %w", err)
				}
				break
			} else {
				return "", fmt.Errorf("error while checking conversation: %w", err)
			}
		}
	}

	// Aggiungi l'utente e l'amico alla conversazione
	_, err = tx.Exec("INSERT INTO members (userId, conversationId) VALUES (?, ?)", id, conversationId)
	if err != nil {
		return "", fmt.Errorf("error while adding user to conversation: %w", err)
	}

	_, err = tx.Exec("INSERT INTO members (userId, conversationId) VALUES (?, ?)", friendId, conversationId)
	if err != nil {
		return "", fmt.Errorf("error while adding friend to conversation: %w", err)
	}

	// Se tutto è andato a buon fine, esegui il commit
	err = db.CloseTx(tx, true)
	if err != nil {
		return "", err
	}
	// Restituisci l'ID della conversazione
	return conversationId, nil
}

func (db *appdbimpl) StartConversationGroup(id string, groupName string) (string, error) {
	// Inizia la transazione
	tx, err := db.GetTx()
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			// Se un errore è presente, esegui il rollback
			rollbackErr := db.CloseTx(tx, false)
			if rollbackErr != nil {
				// In caso di errore durante il rollback, registriamo un errore
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
		}
	}()

	// Genera un nuovo ID per la conversazione
	var conversationId string
	for {
		conversationId = uuid.New().String()
		err := tx.QueryRow("SELECT id FROM conversations WHERE id = ?", conversationId).Scan(&conversationId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Inserisci una nuova conversazione nel database
				_, err = tx.Exec("INSERT INTO conversations (id, type, name) VALUES (?, 'group', ?)", conversationId, groupName)
				if err != nil {
					return "", fmt.Errorf("error while inserting conversation: %w", err)
				}
				break
			} else {
				return "", fmt.Errorf("error while checking conversation: %w", err)
			}
		}
	}

	// Aggiungi l'utente alla conversazione
	_, err = tx.Exec("INSERT INTO members (userId, conversationId) VALUES (?, ?)", id, conversationId)
	if err != nil {
		return "", fmt.Errorf("error while adding user to conversation: %w", err)
	}

	// Se tutte le operazioni sono andate a buon fine, esegui il commit
	err = db.CloseTx(tx, true)
	if err != nil {
		return "", err
	}

	// Restituisci l'ID della conversazione
	return conversationId, nil
}
