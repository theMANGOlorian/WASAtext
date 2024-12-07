package database

import (
	"database/sql"
	"fmt"
)

// GetTx restituisce una nuova transazione.
func (db *appdbimpl) GetTx() (*sql.Tx, error) {
	// Avvia una nuova transazione con db.c, che è il database connection pool
	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	return tx, nil
}

func (db *appdbimpl) CloseTx(tx *sql.Tx, success bool) error {
	// Se success è true, effettua il commit della transazione
	if success {
		return tx.Commit()
	}
	// Se success è false, effettua il rollback della transazione
	return tx.Rollback()
}
