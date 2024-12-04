/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"WASAtext/service/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// prof example
	GetName() (string, error)
	SetName(name string) error

	/* My interfaces */                                       // Remember: use Capital letter -> to set visibility
	DoLogin(username string) (string, error)                  // trying to log in
	SetMyUserName(id string, username string) (string, error) // change username
	SetMyPhoto(id string) (string, error)                     // Change photo profile
	StartConversationPrivate(id string, friendName string) (string, error)
	StartConversationGroup(id string, groupName string) (string, error)
	GetConversations(id string, r *utils.GetConversationsResponseBody) error
	AddToGroupPermission(userId string, groupId string) (int, error)
	AddToGroup(friendId string, groupId string) (int, error)
	LeaveGroup(userId string, groupId string) (int, error)
	SetGroupName(userId string, groupId string, name string) (int, error)
	SetGroupPhoto(userId string, groupId string) (string, int, error)
	SendMessage(userId string, conversationId string, text string, replyTo ...string) (int, *utils.SendMessageResponseBody, error)
	// special interface
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB, schemaFilePath string) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error checking table existence: %w", err)
	}

	// If the tables don't exist, apply migrations from the provided SQL file
	if err == sql.ErrNoRows {
		// Read the SQL file
		sqlBytes, err := ioutil.ReadFile(schemaFilePath)
		if err != nil {
			return nil, fmt.Errorf("error reading SQL file: %w", err)
		}

		// Execute SQL file content
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return nil, fmt.Errorf("error applying SQL migrations: %w", err)
		}
		log.Println("Database tables initialized")
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
